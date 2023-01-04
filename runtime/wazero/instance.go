/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wazero

import (
	"context"
	"sync"
	"sync/atomic"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	klog "k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/abi"
	pwapi "github.com/banzaicloud/proxy-wasm-go-host/api"
	"github.com/banzaicloud/proxy-wasm-go-host/internal/imports"
)

var (
	ErrInstanceNotStart       = errors.New("instance has not started")
	ErrInstanceAlreadyStart   = errors.New("instance has already started")
	ErrOutOfMemory            = errors.New("out of memory")
	ErrUnableToReadMemory     = errors.New("unable to read memory")
	ErrUnknownFunc            = errors.New("unknown func")
	ErrInvalidReturnAddress   = errors.New("invalid return address")
	ErrMallocFunctionNotFound = errors.New("could not find memory allocate function")
)

type Instance struct {
	ctx    context.Context
	vm     *VM
	module *Module

	namespace wazero.Namespace
	instance  api.Module

	lock     sync.Mutex
	started  uint32
	refCount int
	stopCond *sync.Cond

	// user-defined data
	data interface{}

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
	// Here, we initialize an empty namespace as imports are defined prior to start.
	ins := &Instance{
		ctx:       ctx,
		vm:        vm,
		module:    module,
		namespace: vm.runtime.NewNamespace(ctx),
		lock:      sync.Mutex{},

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

func (i *Instance) GetModule() pwapi.WasmModule {
	return i.module
}

// Start makes a new namespace which has the module dependencies of the guest.
func (i *Instance) Start() error {
	if i.checkStart() {
		return ErrInstanceAlreadyStart
	}

	if err := i.registerImports(); err != nil {
		return err
	}

	ctx := context.Background()
	r := i.vm.runtime
	ns := i.namespace

	if _, err := wasi_snapshot_preview1.NewBuilder(r).Instantiate(ctx, ns); err != nil {
		ns.Close(ctx)
		i.logger.Error(err, "could not instantiate wasi_snapshot_preview1")
		return err
	}

	ins, err := ns.InstantiateModule(ctx, i.module.module, wazero.NewModuleConfig())
	if err != nil {
		ns.Close(ctx)
		i.logger.Error(err, "could not instantiate module")
		return err
	}

	i.instance = ins

	for _, fn := range i.startFunctionNames {
		f := i.instance.ExportedFunction(fn)
		if f == nil {
			continue
		}

		if _, err := f.Call(context.Background()); err != nil {
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
		atomic.CompareAndSwapUint32(&i.started, 1, 0)
		i.lock.Unlock()

		if ns := i.namespace; ns != nil {
			ns.Close(i.ctx)
		}
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

	r := i.vm.runtime
	ns := i.namespace

	// proxy-wasm cannot run multiple ABI in the same instance because the ABI
	// collides. They all use the same module name: "env"
	module := "env"

	var hostFunctions func(pwapi.WasmInstance) map[string]interface{}

	abiName := abi.ProxyWasmABI_0_2_1
	if abiList := i.module.GetABINameList(); len(abiList) > 0 {
		abiName = abiList[0]
	}

	// Instantiate WASI also under the unstable name for old compilers,
	// such as TinyGo 0.19 used for v1 ABI.
	if abiName == abi.ProxyWasmABI_0_1_0 {
		wasiBuilder := r.NewHostModuleBuilder("wasi_unstable")
		wasi_snapshot_preview1.NewFunctionExporter().ExportFunctions(wasiBuilder)
		if _, err := wasiBuilder.Instantiate(i.ctx, ns); err != nil {
			ns.Close(i.ctx)
			i.logger.Error(err, "could not instantiate wasi_unstable")
			return err
		}
	}

	hostFunctions = imports.HostFunctions
	b := r.NewHostModuleBuilder(module)
	for n, f := range hostFunctions(i) {
		b.NewFunctionBuilder().WithFunc(f).Export(n)
	}

	if _, err := b.Instantiate(i.ctx, ns); err != nil {
		i.logger.Error(err, "could not instantiate module")
		return err
	}

	return nil
}

func (i *Instance) Malloc(size int32) (uint64, error) {
	if !i.checkStart() {
		return 0, ErrInstanceNotStart
	}

	var f api.Function
	mallocFuncNames := i.mallocFunctionNames
	for _, fn := range mallocFuncNames {
		if f == nil {
			f = i.instance.ExportedFunction(fn)
		}
		if f != nil {
			break
		}
	}

	if f == nil {
		return 0, ErrMallocFunctionNotFound
	}

	addr, err := i.GetWasmFunction(f).Call(size)
	if err != nil {
		i.HandleError(err)
		return 0, err
	}

	if v, ok := addr.(int32); ok {
		return uint64(v), nil
	}

	return 0, ErrInvalidReturnAddress
}

func (i *Instance) GetExportsFunc(funcName string) (pwapi.WasmFunction, error) {
	if !i.checkStart() {
		return nil, ErrInstanceNotStart
	}

	wf := i.instance.ExportedFunction(funcName)
	if wf == nil {
		return nil, ErrUnknownFunc
	}

	return i.GetWasmFunction(wf), nil
}

func (i *Instance) GetWasmFunction(f api.Function) pwapi.WasmFunction {
	fn := &wasmFunction{
		fn:     f,
		logger: i.logger,
	}

	if rts := f.Definition().ResultTypes(); len(rts) > 0 {
		fn.rt = rts[0]
	}

	return fn
}

type wasmFunction struct {
	fn     api.Function
	rt     api.ValueType
	logger logr.Logger
}

// Call implements api.WasmFunction
func (f *wasmFunction) Call(args ...interface{}) (interface{}, error) {
	realArgs := make([]uint64, 0, len(args))
	for _, a := range args {
		if _, v, err := convertFromGoValue(a); err != nil {
			return nil, err
		} else {
			realArgs = append(realArgs, v)
		}
	}

	if len(f.fn.Definition().ExportNames()) > 0 {
		f.logger.V(3).Info("call module function", "name", f.fn.Definition().ExportNames()[0], "args", realArgs)
	}

	if ret, err := f.fn.Call(context.Background(), realArgs...); err != nil {
		return nil, err
	} else if len(ret) == 0 {
		return nil, nil
	} else {
		return convertToGoValue(f.rt, ret[0])
	}
}

func (i *Instance) GetExportsMem(memName string) ([]byte, error) {
	if !i.checkStart() {
		return nil, ErrInstanceNotStart
	}

	mem := i.instance.ExportedMemory(memName)

	return i.GetMemory(0, uint64(mem.Size(i.ctx)))
}

func (i *Instance) GetMemory(addr uint64, size uint64) ([]byte, error) {
	if ret, ok := i.instance.Memory().Read(i.ctx, uint32(addr), uint32(size)); !ok {
		return nil, ErrUnableToReadMemory
	} else {
		return ret, nil
	}
}

func (i *Instance) PutMemory(addr uint64, size uint64, content []byte) error {
	if n := len(content); n < int(size) {
		size = uint64(n)
	}

	if ok := i.instance.Memory().Write(i.ctx, uint32(addr), content[:size]); !ok {
		return ErrOutOfMemory
	}

	return nil
}

func (i *Instance) GetByte(addr uint64) (byte, error) {
	if b, ok := i.instance.Memory().ReadByte(i.ctx, uint32(addr)); !ok {
		return b, ErrOutOfMemory
	} else {
		return b, nil
	}
}

func (i *Instance) PutByte(addr uint64, b byte) error {
	if ok := i.instance.Memory().WriteByte(i.ctx, uint32(addr), b); !ok {
		return ErrOutOfMemory
	}

	return nil
}

func (i *Instance) GetUint32(addr uint64) (uint32, error) {
	if n, ok := i.instance.Memory().ReadUint32Le(i.ctx, uint32(addr)); !ok {
		return n, ErrOutOfMemory
	} else {
		return n, nil
	}
}

func (i *Instance) PutUint32(addr uint64, value uint32) error {
	if ok := i.instance.Memory().WriteUint32Le(i.ctx, uint32(addr), value); !ok {
		return ErrOutOfMemory
	}

	return nil
}

func (i *Instance) HandleError(err error) {
	i.logger.Error(err, "wasm error")
}
