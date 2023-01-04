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

package e2e

import (
	"context"
	_ "embed"
	"testing"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
	"github.com/banzaicloud/proxy-wasm-go-host/runtime/wazero"
)

func BenchmarkStartABIContext_wazero(b *testing.B) {
	vm := wazero.NewVM(context.Background())
	defer vm.Close()

	benchmarkStartABIContext(b, vm)
}

func benchmarkStartABIContext(b *testing.B, vm api.WasmVM) {
	b.Helper()

	module, err := vm.NewModule(binAddRequestHeader)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		instance, err := module.NewInstance()
		if err != nil {
			b.Fatal(err)
		}

		if _, err := startABIContext(instance); err != nil {
			b.Fatal(err)
		} else {
			instance.Stop()
		}
	}
}

func BenchmarkAddRequestHeader_wazero(b *testing.B) {
	vm := wazero.NewVM(context.Background())
	defer vm.Close()

	benchmarkAddRequestHeader(b, vm)
}

func benchmarkAddRequestHeader(b *testing.B, vm api.WasmVM) {
	b.Helper()

	module, err := vm.NewModule(binAddRequestHeader)
	if err != nil {
		b.Fatal(err)
	}

	instance, err := module.NewInstance()
	if err != nil {
		b.Fatal(err)
	}

	defer instance.Stop()

	benchmark(b, instance, testAddRequestHeader)
}

func benchmark(b *testing.B, instance api.WasmInstance, test func(wasmCtx api.ABIContext, contextID int32) error) {
	b.Helper()

	wasmCtx, err := startABIContext(instance)
	if err != nil {
		b.Fatal(err)
	}
	defer wasmCtx.GetInstance().Stop()

	exports := wasmCtx.GetExports()

	// make the root context
	rootContextID := int32(1)
	if err = exports.ProxyOnContextCreate(rootContextID, int32(0)); err != nil {
		b.Fatal(err)
	}

	// lock wasm vm instance for exclusive ownership
	wasmCtx.GetInstance().Lock(wasmCtx)
	defer wasmCtx.GetInstance().Unlock()

	// Time the guest call for context create and delete, which happens per-request.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		contextID := int32(2)
		if err = exports.ProxyOnContextCreate(contextID, rootContextID); err != nil {
			b.Fatal(err)
		}

		if err = test(wasmCtx, contextID); err != nil {
			b.Fatal(err)
		}

		if _, err = exports.ProxyOnDone(contextID); err != nil {
			b.Fatal(err)
		}

		if err = exports.ProxyOnDelete(contextID); err != nil {
			b.Fatal(err)
		}
	}
}
