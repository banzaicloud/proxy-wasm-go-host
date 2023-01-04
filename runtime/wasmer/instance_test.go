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

//nolint:forcetypeassert
package wasmer

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
)

func Test_registerFunc(t *testing.T) {
	t.Parallel()

	vm := NewWasmerVM()
	defer vm.Close()

	assert.Equal(t, vm.Name(), "wasmer")

	module, err := vm.NewModule([]byte(`
		(module
			(memory (;0;) 1 1)
			(export "memory" (memory 0))
			(import "wasi_snapshot_preview1" "args_get" (func (param i32 i32) (result i32)))
			(func (export "_start")))
	`))
	require.Nil(t, err)

	_ins, err := module.NewInstance()
	require.Nil(t, err)
	ins := _ins.(*Instance)

	defer ins.Stop()

	// invalid namespace
	assert.Equal(t, ins.registerFunc("", "funcName", nil), ErrInvalidParam)

	// nil f
	assert.Equal(t, ins.registerFunc("TestRegisterFuncNamespace", "funcName", nil), ErrInvalidParam)

	var testStruct struct{}

	// f is not func
	assert.Equal(t, ins.registerFunc("TestRegisterFuncNamespace", "funcName", &testStruct), ErrRegisterNotFunc)

	assert.Nil(t, ins.registerFunc("TestRegisterFuncNamespace", "funcName", func() {}))

	assert.Nil(t, ins.Start())
}

func Test_registerFuncRecoverPanic(t *testing.T) {
	t.Parallel()

	vm := NewWasmerVM()
	defer vm.Close()

	module, err := vm.NewModule([]byte(`
			(module
				(memory (;0;) 1 1)
				(export "memory" (memory 0))
				(import "TestRegisterFuncRecover" "somePanic" (func $somePanic (result i32)))
				(import "wasi_snapshot_preview1" "args_get" (func (param i32 i32) (result i32)))
				(func (export "_start"))
				(func (export "panicTrigger") (param) (result i32)
					call $somePanic))
	`))
	require.Nil(t, err)

	_ins, err := module.NewInstance()
	require.Nil(t, err)
	ins := _ins.(*Instance)

	defer ins.Stop()

	assert.Nil(t, ins.registerFunc("TestRegisterFuncRecover", "somePanic", func() int32 {
		panic("some panic")
	}))

	assert.Nil(t, ins.Start())

	f, err := ins.GetExportsFunc("panicTrigger")
	assert.Nil(t, err)

	_, err = f.Call()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "panic [some panic] when calling func [somePanic]")
}

func TestInstanceMalloc(t *testing.T) {
	t.Parallel()

	vm := NewWasmerVM()
	defer vm.Close()

	module, err := vm.NewModule([]byte(`
			(module
				(memory (;0;) 1 1)
				(export "memory" (memory 0))
				(import "wasi_snapshot_preview1" "args_get" (func (param i32 i32) (result i32)))
				(func (export "_start"))
				(func (export "malloc") (param i32) (result i32) i32.const 10))
	`))
	require.Nil(t, err)

	_ins, err := module.NewInstance()
	require.Nil(t, err)
	ins := _ins.(*Instance)

	defer ins.Stop()

	assert.Nil(t, ins.registerFunc("TestRegisterFuncRecover", "somePanic", func() int32 {
		panic("some panic")
	}))

	assert.Nil(t, ins.Start())

	addr, err := ins.Malloc(100)
	assert.Nil(t, err)
	assert.Equal(t, addr, uint64(10))
}

func TestInstanceMem(t *testing.T) {
	t.Parallel()

	vm := NewWasmerVM()
	defer vm.Close()

	module, err := vm.NewModule([]byte(`
		(module
			(memory (;0;) 1 1)
			(export "memory" (memory 0))
			(import "wasi_snapshot_preview1" "args_get" (func (param i32 i32) (result i32)))
			(func (export "_start")))
		`))
	require.Nil(t, err)

	ins, err := module.NewInstance()
	require.Nil(t, err)

	defer ins.Stop()

	assert.Nil(t, ins.Start())

	m, err := ins.GetExportsMem("memory")
	assert.Nil(t, err)
	// A WebAssembly page has a constant size of 65,536 bytes, i.e., 64KiB
	assert.Equal(t, len(m), 1<<16)

	assert.Nil(t, ins.PutByte(uint64(100), 'a'))
	b, err := ins.GetByte(uint64(100))
	assert.Nil(t, err)
	assert.Equal(t, b, byte('a'))

	assert.Nil(t, ins.PutUint32(uint64(200), 99))
	u, err := ins.GetUint32(uint64(200))
	assert.Nil(t, err)
	assert.Equal(t, u, uint32(99))

	assert.Nil(t, ins.PutMemory(uint64(300), 10, []byte("1111111111")))
	bs, err := ins.GetMemory(uint64(300), 10)
	assert.Nil(t, err)
	assert.Equal(t, string(bs), "1111111111")
}

func TestInstanceData(t *testing.T) {
	t.Parallel()

	vm := NewWasmerVM()
	defer vm.Close()

	module, err := vm.NewModule([]byte(`
			(module
				(memory (;0;) 1 1)
				(export "memory" (memory 0))
				(import "wasi_snapshot_preview1" "args_get" (func (param i32 i32) (result i32)))
				(func (export "_start")))
	`))
	require.Nil(t, err)

	ins, err := module.NewInstance()
	require.Nil(t, err)

	defer ins.Stop()

	assert.Nil(t, ins.Start())

	var data int = 1
	ins.SetData(data)
	assert.Equal(t, ins.GetData().(int), 1)

	for i := 0; i < 10; i++ {
		ins.Lock(i)
		assert.Equal(t, ins.GetData().(int), i)
		ins.Unlock()
	}
}

func TestWasmerTypes(t *testing.T) {
	t.Parallel()

	testDatas := []struct {
		refType     reflect.Type
		refValue    reflect.Value
		refValKind  reflect.Kind
		wasmValKind wasmerGo.ValueKind
	}{
		{reflect.TypeOf(int32(0)), reflect.ValueOf(int32(0)), reflect.Int32, wasmerGo.I32},
		{reflect.TypeOf(int64(0)), reflect.ValueOf(int64(0)), reflect.Int64, wasmerGo.I64},
		{reflect.TypeOf(float32(0)), reflect.ValueOf(float32(0)), reflect.Float32, wasmerGo.F32},
		{reflect.TypeOf(float64(0)), reflect.ValueOf(float64(0)), reflect.Float64, wasmerGo.F64},
	}

	for _, tc := range testDatas {
		assert.Equal(t, convertFromGoType(tc.refType).Kind(), tc.wasmValKind)
		assert.Equal(t, convertToGoTypes(convertFromGoValue(tc.refValue)).Kind(), tc.refValKind)
	}
}

func TestRefCount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vm := NewWasmerVM()
	defer vm.Close()

	module, err := vm.NewModule([]byte(`
		(module
			(memory (;0;) 1 1)
			(export "memory" (memory 0))
			(import "wasi_snapshot_preview1" "args_get" (func (param i32 i32) (result i32)))
			(func (export "_start")))
	`))
	require.Nil(t, err)

	ins, err := NewWasmerInstance(vm.(*VM), module.(*Module))
	require.Nil(t, err)

	assert.False(t, ins.Acquire())

	ins.started = 1
	for i := 0; i < 100; i++ {
		assert.True(t, ins.Acquire())
	}
	assert.Equal(t, ins.refCount, 100)

	ins.Stop()
	ins.Stop() // double stop
	time.Sleep(time.Second)
	assert.Equal(t, ins.started, uint32(1))

	for i := 0; i < 100; i++ {
		ins.Release()
	}

	time.Sleep(time.Second)
	assert.False(t, ins.Acquire())
	assert.Equal(t, ins.started, uint32(0))
	assert.Equal(t, ins.refCount, 0)
}
