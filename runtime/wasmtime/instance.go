// Copyright (c) 2023 Cisco and/or its affiliates. All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package wasmtime

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"

	"emperror.dev/errors"
	"github.com/bytecodealliance/wasmtime-go/v3"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
	"github.com/banzaicloud/proxy-wasm-go-host/internal/imports"
)

var (
	ErrAddrOverflow         = errors.New("addr overflow")
	ErrInstanceNotStart     = errors.New("instance has not started")
	ErrInstanceAlreadyStart = errors.New("instance has already started")
	ErrInvalidParam         = errors.New("invalid param")
	ErrRegisterNotFunc      = errors.New("register a non-func object")
	ErrRegisterArgType      = errors.New("register func with invalid arg type")

	ErrPutUint32                   = errors.New("could not put uint32 to memory")
	ErrGetUint32                   = errors.New("could not get uint32 from memory")
	ErrPutByte                     = errors.New("could not put byte to memory")
	ErrGetByte                     = errors.New("could not get byte from memory")
	ErrPutData                     = errors.New("could not put data to memory")
	ErrGetData                     = errors.New("could not get data to memory")
	ErrGetExportsMemNotImplemented = errors.New("GetExportsMem not implemented")
	ErrMallocFunctionNotFound      = errors.New("could not find memory allocate function")
	ErrMalloc                      = errors.New("could not allocate memory")
)

type Instance struct {
	lock     sync.Mutex
	started  uint32
	refCount int
	stopCond *sync.Cond

	// user-defined data
	data interface{}

	store    *wasmtime.Store
	module   *Module
	instance *wasmtime.Instance
	memory   *wasmtime.Memory

	hostModules map[string][]hostFunc

	logger              logr.Logger
	startFunctionNames  []string
	mallocFunctionNames []string
}

type InstanceOptions func(instance *Instance)

func InstanceWithStartFunctionNames(names ...string) InstanceOptions {
	return func(instance *Instance) {
		instance.startFunctionNames = names
	}
}

func InstanceWithMallocFunctionNames(names ...string) InstanceOptions {
	return func(instance *Instance) {
		instance.mallocFunctionNames = names
	}
}

func InstanceWithLogger(logger logr.Logger) InstanceOptions {
	return func(instance *Instance) {
		instance.logger = logger
	}
}

func NewInstance(ctx context.Context, vm *VM, module *Module, options ...InstanceOptions) *Instance {
	ins := &Instance{
		store:               wasmtime.NewStore(vm.engine),
		module:              module,
		lock:                sync.Mutex{},
		hostModules:         make(map[string][]hostFunc),
		startFunctionNames:  []string{"_start", "_initialize"},
		mallocFunctionNames: []string{"proxy_on_memory_allocate", "malloc"},
	}
	ins.stopCond = sync.NewCond(&ins.lock)

	for _, option := range options {
		option(ins)
	}

	if ins.logger == (logr.Logger{}) {
		ins.logger = klog.Background()
	}

	return ins
}

func (i *Instance) GetData() interface{} {
	return i.data
}

func (i *Instance) SetData(data interface{}) {
	i.data = data
}

func (i *Instance) Acquire() bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	if !i.checkStart() {
		return false
	}

	i.refCount++

	return true
}

func (i *Instance) Release() {
	i.lock.Lock()
	if i.refCount > 0 {
		i.refCount--
	}

	if i.refCount <= 0 {
		i.stopCond.Broadcast()
	}
	i.lock.Unlock()
}

func (i *Instance) Lock(data interface{}) {
	i.lock.Lock()
	i.data = data
}

func (i *Instance) Unlock() {
	i.data = nil
	i.lock.Unlock()
}

func (i *Instance) GetModule() api.WasmModule {
	return i.module
}

func (i *Instance) Start() error {
	if i.checkStart() {
		return ErrInstanceAlreadyStart
	}

	if err := i.registerImports(); err != nil {
		return err
	}

	for _, fn := range i.startFunctionNames {
		f := i.instance.GetFunc(i.store, fn)
		if f == nil {
			continue
		}

		if _, err := f.Call(i.store); err != nil {
			i.HandleError(err)
			return err
		}

		atomic.StoreUint32(&i.started, 1)

		return nil
	}

	return errors.NewWithDetails("could not start instance: start function is not exported", "functions", i.startFunctionNames)
}

func (i *Instance) Stop() {
	go func() {
		i.lock.Lock()
		for i.refCount > 0 {
			i.stopCond.Wait()
		}
		_ = atomic.CompareAndSwapUint32(&i.started, 1, 0)
		// if err := i.instance.Close(i.moduleCtx); err != nil { // TODO?
		// 	i.Logger().Error(err, "could not close module", "module", i.module.Name())
		// }
		i.lock.Unlock()
	}()
}

// return true is Instance is started, false if not started.
func (i *Instance) checkStart() bool {
	return atomic.LoadUint32(&i.started) == 1
}

func (i *Instance) registerImports() error {
	if i.checkStart() {
		return ErrInstanceAlreadyStart
	}

	linker := wasmtime.NewLinker(i.store.Engine)

	err := linker.DefineWasi()
	if err != nil {
		return err
	}

	i.store.SetWasi(wasmtime.NewWasiConfig())

	// proxy-wasm cannot run multiple ABI in the same instance because the ABI
	// collides. They all use the same module name: "env"
	module := "env"

	hostFunctions := imports.HostFunctions
	for n, f := range hostFunctions(i) {
		err := linker.DefineFunc(i.store, module, n, f)
		if err != nil {
			return err
		}
	}

	i.instance, err = linker.Instantiate(i.store, i.module.module)
	if err != nil {
		return err
	}

	i.memory = i.instance.GetExport(i.store, "memory").Memory()

	return nil
}

type hostFunc struct {
	funcName string
	f        interface{}
}

func (i *Instance) RegisterFunc(namespace string, funcName string, f interface{}) error {
	if i.checkStart() {
		return ErrInstanceAlreadyStart
	}

	if _, ok := i.hostModules[namespace]; !ok {
		i.hostModules[namespace] = []hostFunc{}
	}

	i.hostModules[namespace] = append(i.hostModules[namespace], hostFunc{funcName: funcName, f: f})

	return nil
}

func (i *Instance) Malloc(size int32) (uint64, error) {
	if !i.checkStart() {
		return 0, ErrInstanceNotStart
	}

	var f *wasmtime.Func
	mallocFuncNames := i.mallocFunctionNames
	var mallocFuncName string
	for _, fn := range mallocFuncNames {
		if f == nil {
			f = i.instance.GetFunc(i.store, fn)
		}
		if f != nil {
			mallocFuncName = fn
			break
		}
	}

	if f == nil {
		return 0, ErrMallocFunctionNotFound
	}

	malloc := &Call{
		Func:   f,
		name:   mallocFuncName,
		store:  i.store,
		logger: i.logger,
	}

	addr, err := malloc.Call(size)
	if err != nil {
		i.HandleError(err)
		return 0, err
	}

	if v, ok := addr.(int32); ok {
		return uint64(v), nil
	}

	return 0, ErrMalloc
}

type Call struct {
	*wasmtime.Func

	name string

	store *wasmtime.Store

	logger logr.Logger
}

func (c *Call) Call(args ...interface{}) (interface{}, error) {
	c.logger.V(3).Info(fmt.Sprintf("call module function %s", c.name))

	ret, err := c.Func.Call(c.store, args...)
	if err != nil {
		return 0, err
	}

	switch r := ret.(type) {
	case []wasmtime.Val:
		return r[0].I32(), nil
	default:
		return r, nil
	}
}

func (i *Instance) GetExportsFunc(funcName string) (api.WasmFunction, error) {
	if !i.checkStart() {
		return nil, ErrInstanceNotStart
	}

	f := &Call{
		Func:   i.instance.GetFunc(i.store, funcName),
		name:   funcName,
		store:  i.store,
		logger: i.logger,
	}

	return f, nil
}

func (i *Instance) GetExportsMem(memName string) ([]byte, error) {
	return i.memory.UnsafeData(i.store), nil
}

func (i *Instance) GetMemory(addr uint64, size uint64) ([]byte, error) {
	memory := i.memory.UnsafeData(i.store)

	if i.memory.DataSize(i.store) <= uintptr(addr+size) {
		return nil, ErrGetData
	}

	return memory[uint32(addr) : uint32(addr)+uint32(size)], nil
}

func (i *Instance) PutMemory(addr uint64, size uint64, content []byte) error {
	if n := len(content); n < int(size) {
		size = uint64(n)
	}

	if i.memory.DataSize(i.store) <= uintptr(addr+size) {
		return ErrPutData
	}

	memory := i.memory.UnsafeData(i.store)

	copy(memory[uint32(addr):], content[:size])

	return nil
}

func (i *Instance) GetByte(addr uint64) (byte, error) {
	bytes, err := i.GetMemory(addr, 1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func (i *Instance) PutByte(addr uint64, b byte) error {
	return i.PutMemory(addr, 1, []byte{b})
}

func (i *Instance) GetUint32(addr uint64) (uint32, error) {
	data, err := i.GetMemory(addr, 4)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(data), nil
}

func (i *Instance) PutUint32(addr uint64, value uint32) error {
	data, err := i.GetMemory(addr, 4)
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(data, value)

	return nil
}

func (i *Instance) HandleError(err error) {
	i.logger.Error(err, "wasm error")
}
