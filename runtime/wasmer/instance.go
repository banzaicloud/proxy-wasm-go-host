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

//nolint:goerr113
package wasmer

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
	"k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/abi"
	"github.com/banzaicloud/proxy-wasm-go-host/api"
	"github.com/banzaicloud/proxy-wasm-go-host/internal/imports"
)

var (
	ErrAddrOverflow           = errors.New("addr overflow")
	ErrInstanceNotStart       = errors.New("instance has not started")
	ErrInstanceAlreadyStart   = errors.New("instance has already started")
	ErrInvalidParam           = errors.New("invalid param")
	ErrRegisterNotFunc        = errors.New("register a non-func object")
	ErrRegisterArgType        = errors.New("register func with invalid arg type")
	ErrInvalidReturnAddress   = errors.New("invalid return address")
	ErrMallocFunctionNotFound = errors.New("could not find memory allocate function")
)

type Instance struct {
	vm           *VM
	module       *Module
	importObject *wasmerGo.ImportObject
	instance     *wasmerGo.Instance
	debug        *dwarfInfo

	lock     sync.Mutex
	started  uint32
	refCount int
	stopCond *sync.Cond

	// for cache
	memory    *wasmerGo.Memory
	funcCache sync.Map // string -> *wasmerGo.Function

	// user-defined data
	data interface{}

	logger              logr.Logger
	startFunctionNames  []string
	mallocFunctionNames []string
}

type InstanceOptions func(instance *Instance)

func InstanceWithDebug(debug *dwarfInfo) InstanceOptions {
	return func(instance *Instance) {
		if debug != nil {
			instance.debug = debug
		}
	}
}

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

func NewWasmerInstance(vm *VM, module *Module, options ...InstanceOptions) (*Instance, error) {
	ins := &Instance{
		vm:     vm,
		module: module,
		lock:   sync.Mutex{},

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

	wasiEnv, err := wasmerGo.NewWasiStateBuilder("").Finalize()
	if err != nil || wasiEnv == nil {
		return nil, errors.WrapIf(err, "could not create wasi env")
	}

	imo, err := wasiEnv.GenerateImportObject(ins.vm.store, ins.module.module)
	if err != nil {
		return nil, errors.WrapIf(err, "could not create import object")
	}

	ins.importObject = imo

	return ins, nil
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
	i.refCount--

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

	ins, err := wasmerGo.NewInstance(i.module.module, i.importObject)
	if err != nil {
		return errors.WrapIf(err, "could not start instance")
	}

	i.instance = ins

	for _, fn := range i.startFunctionNames {
		f, err := i.instance.Exports.GetFunction(fn)
		if err != nil {
			continue
		}

		if _, err := f(); err != nil {
			i.HandleError(err)
			return errors.WrapIf(err, "could not call start function")
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

		atomic.CompareAndSwapUint32(&i.started, 1, 0)
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

	// proxy-wasm cannot run multiple ABI in the same instance because the ABI
	// collides. They all use the same module name: "env"
	module := "env"

	abiName := abi.ProxyWasmABI_0_2_1
	if abiList := i.module.GetABINameList(); len(abiList) > 0 {
		abiName = abiList[0]
	}

	var hostFunctions func(api.WasmInstance) map[string]interface{}
	switch abiName {
	case abi.ProxyWasmABI_0_2_1, abi.ProxyWasmABI_0_2_0, abi.ProxyWasmABI_0_1_0:
		hostFunctions = imports.HostFunctions
	default:
		return fmt.Errorf("unknown ABI: %s", abiName)
	}

	for n, f := range hostFunctions(i) {
		if err := i.registerFunc(module, n, f); err != nil {
			return errors.WrapIfWithDetails(err, "could not register func", "name", n, "namespace", module, "func", f)
		}
	}

	return nil
}

func (i *Instance) registerFunc(namespace string, funcName string, f interface{}) error {
	if namespace == "" || funcName == "" {
		return ErrInvalidParam
	}

	if f == nil || reflect.ValueOf(f).IsNil() {
		return ErrInvalidParam
	}

	if reflect.TypeOf(f).Kind() != reflect.Func {
		return ErrRegisterNotFunc
	}

	funcType := reflect.TypeOf(f)

	argsNum := funcType.NumIn()

	argsKind := make([]*wasmerGo.ValueType, argsNum)
	for i := 0; i < argsNum; i++ {
		argsKind[i] = convertFromGoType(funcType.In(i))
	}

	retsNum := funcType.NumOut()
	retsKind := make([]*wasmerGo.ValueType, retsNum)
	for i := 0; i < retsNum; i++ {
		retsKind[i] = convertFromGoType(funcType.Out(i))
	}

	fwasmer := wasmerGo.NewFunction(
		i.vm.store,
		wasmerGo.NewFunctionType(argsKind, retsKind),
		func(args []wasmerGo.Value) (callRes []wasmerGo.Value, err error) {
			defer func() {
				if r := recover(); r != nil {
					callRes = nil
					err = fmt.Errorf("panic [%v] when calling func [%v]", r, funcName)
				}
			}()

			callArgs := make([]reflect.Value, len(args))

			for i, arg := range args {
				callArgs[i] = convertToGoTypes(arg)
			}

			callResult := reflect.ValueOf(f).Call(callArgs)

			ret := convertFromGoValue(callResult[0])

			return []wasmerGo.Value{ret}, err
		},
	)

	i.importObject.Register(namespace, map[string]wasmerGo.IntoExtern{
		funcName: fwasmer,
	})

	return nil
}

func (i *Instance) Malloc(size int32) (uint64, error) {
	if !i.checkStart() {
		return 0, ErrInstanceNotStart
	}

	var f api.WasmFunction
	mallocFuncNames := i.mallocFunctionNames
	for _, fn := range mallocFuncNames {
		if fn, err := i.GetExportsFunc(fn); err == nil {
			f = fn
			break
		}
	}

	if f == nil {
		return 0, ErrMallocFunctionNotFound
	}

	addr, err := f.Call(size)
	if err != nil {
		i.HandleError(err)
		return 0, err
	}

	if v, ok := addr.(int32); ok {
		return uint64(v), nil
	}

	return 0, ErrInvalidReturnAddress
}

func (i *Instance) GetExportsFunc(funcName string) (api.WasmFunction, error) {
	if !i.checkStart() {
		return nil, ErrInstanceNotStart
	}

	if v, ok := i.funcCache.Load(funcName); ok {
		if f, ok := v.(*wasmerGo.Function); ok {
			return &wasmFunction{name: funcName, logger: i.logger, fn: f}, nil
		}
	}

	f, err := i.instance.Exports.GetRawFunction(funcName)
	if err != nil {
		return nil, err
	}

	i.funcCache.Store(funcName, f)

	return &wasmFunction{name: funcName, logger: i.logger, fn: f}, nil
}

type wasmFunction struct {
	fn     *wasmerGo.Function
	name   string
	logger logr.Logger
}

func (f *wasmFunction) Call(args ...interface{}) (interface{}, error) {
	f.logger.V(3).Info("call module function", "name", f.name, "args", args)

	return f.fn.Call(args...)
}

func (i *Instance) GetExportsMem(memName string) ([]byte, error) {
	if !i.checkStart() {
		return nil, ErrInstanceNotStart
	}

	if i.memory == nil {
		m, err := i.instance.Exports.GetMemory(memName)
		if err != nil {
			return nil, err
		}

		i.memory = m
	}

	return i.memory.Data(), nil
}

func (i *Instance) GetMemory(addr uint64, size uint64) ([]byte, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return nil, err
	}

	if int(addr) > len(mem) || int(addr+size) > len(mem) {
		return nil, ErrAddrOverflow
	}

	return mem[addr : addr+size], nil
}

func (i *Instance) PutMemory(addr uint64, size uint64, content []byte) error {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if int(addr) > len(mem) || int(addr+size) > len(mem) {
		return ErrAddrOverflow
	}

	copySize := uint64(len(content))
	if size < copySize {
		copySize = size
	}

	copy(mem[addr:], content[:copySize])

	return nil
}

func (i *Instance) GetByte(addr uint64) (byte, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}

	if int(addr) > len(mem) {
		return 0, ErrAddrOverflow
	}

	return mem[addr], nil
}

func (i *Instance) PutByte(addr uint64, b byte) error {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if int(addr) > len(mem) {
		return ErrAddrOverflow
	}

	mem[addr] = b

	return nil
}

func (i *Instance) GetUint32(addr uint64) (uint32, error) {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return 0, err
	}

	if int(addr) > len(mem) || int(addr+4) > len(mem) {
		return 0, ErrAddrOverflow
	}

	return binary.LittleEndian.Uint32(mem[addr:]), nil
}

func (i *Instance) PutUint32(addr uint64, value uint32) error {
	mem, err := i.GetExportsMem("memory")
	if err != nil {
		return err
	}

	if int(addr) > len(mem) || int(addr+4) > len(mem) {
		return ErrAddrOverflow
	}

	binary.LittleEndian.PutUint32(mem[addr:], value)

	return nil
}

func (i *Instance) HandleError(err error) {
	var trapError *wasmerGo.TrapError
	if !errors.As(err, &trapError) {
		return
	}

	trace := trapError.Trace()
	if trace == nil {
		return
	}

	i.logger.Error(err, "wasm error")

	var traceOutput string

	if i.debug == nil {
		// do not have dwarf debug info
		for _, t := range trace {
			traceOutput += fmt.Sprintf("funcIndex: %v, funcOffset: 0x%08x, moduleOffset: 0x%08x",
				t.FunctionIndex(), t.FunctionOffset(), t.ModuleOffset())
		}
	} else {
		for _, t := range trace {
			pc := uint64(t.ModuleOffset())
			line := i.debug.SeekPC(pc)
			if line != nil {
				traceOutput += fmt.Sprintf("funcIndex: %v, funcOffset: 0x%08x, pc: 0x%08x %v:%v",
					t.FunctionIndex(), t.FunctionOffset(), pc, line.File.Name, line.Line)
			} else {
				traceOutput += fmt.Sprintf("funcIndex: %v, funcOffset: 0x%08x, pc: 0x%08x fail to seek pc",
					t.FunctionIndex(), t.FunctionOffset(), t.ModuleOffset())
			}
		}
	}
}
