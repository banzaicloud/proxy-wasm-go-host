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
package wasmtime

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/wabin/binary"
	"github.com/tetratelabs/wabin/wasm"
)

var simpleWasm = binary.EncodeModule(&wasm.Module{
	TypeSection:     []*wasm.FunctionType{{}},                     // v_v
	FunctionSection: []wasm.Index{wasm.Index(0)},                  // type[0] == v_v
	CodeSection:     []*wasm.Code{{Body: []byte{wasm.OpcodeEnd}}}, // noop
	MemorySection:   &wasm.Memory{Min: 1, Max: 1, IsMaxEncoded: true},
	ExportSection: []*wasm.Export{
		{Name: "_start", Type: wasm.ExternTypeFunc, Index: wasm.Index(0)},   // export func[0]
		{Name: "memory", Type: wasm.ExternTypeMemory, Index: wasm.Index(0)}, // export memory[0]
	},
})

func TestInstanceMem(t *testing.T) {
	t.Parallel()

	vm := NewVM(context.Background())
	defer vm.Close()

	module, err := vm.NewModule(simpleWasm)
	require.Nil(t, err)

	ins, err := module.NewInstance()
	require.Nil(t, err)

	defer ins.Stop()

	require.Nil(t, ins.Start())

	m, err := ins.GetExportsMem("memory")
	require.Nil(t, err)
	// A WebAssembly page has a constant size of 65,536 bytes, i.e., 64KiB
	require.Equal(t, len(m), 1<<16)

	require.Nil(t, ins.PutByte(uint64(100), 'a'))
	b, err := ins.GetByte(uint64(100))
	require.Nil(t, err)
	require.Equal(t, b, byte('a'))

	require.Nil(t, ins.PutUint32(uint64(200), 99))
	u, err := ins.GetUint32(uint64(200))
	require.Nil(t, err)
	require.Equal(t, u, uint32(99))

	require.Nil(t, ins.PutMemory(uint64(300), 10, []byte("1111111111")))
	bs, err := ins.GetMemory(uint64(300), 10)
	require.Nil(t, err)
	require.Equal(t, string(bs), "1111111111")
}

func TestInstanceData(t *testing.T) {
	t.Parallel()

	vm := NewVM(context.Background())
	defer vm.Close()

	module, err := vm.NewModule(simpleWasm)
	require.Nil(t, err)

	ins, err := module.NewInstance()
	require.Nil(t, err)

	defer ins.Stop()

	require.Nil(t, ins.Start())

	var data int = 1
	ins.SetData(data)
	require.Equal(t, ins.GetData().(int), 1)

	for i := 0; i < 10; i++ {
		ins.Lock(i)
		require.Equal(t, ins.GetData().(int), i)
		ins.Unlock()
	}
}

func TestRefCount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vm := NewVM(context.Background())
	defer vm.Close()

	module, err := vm.NewModule(simpleWasm)
	require.Nil(t, err)

	ins := NewInstance(context.Background(), vm.(*VM), module.(*Module))

	require.False(t, ins.Acquire())

	ins.started = 1
	for i := 0; i < 100; i++ {
		require.True(t, ins.Acquire())
	}
	require.Equal(t, ins.refCount, 100)

	ins.Stop()
	ins.Stop() // double stop
	time.Sleep(time.Second)
	require.Equal(t, ins.started, uint32(1))

	for i := 0; i < 100; i++ {
		ins.Release()
	}

	time.Sleep(time.Second)
	require.False(t, ins.Acquire())
	require.Equal(t, ins.started, uint32(0))
	require.Equal(t, ins.refCount, 0)
}
