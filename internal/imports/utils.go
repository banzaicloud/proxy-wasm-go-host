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

package imports

import (
	"github.com/banzaicloud/proxy-wasm-go-host/abi"
	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

func copyIntoInstance(instance api.WasmInstance, value string, retPtr int32, retSize int32) api.WasmResult {
	addr, err := instance.Malloc(int32(len(value)))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess
	}

	err = instance.PutMemory(addr, uint64(len(value)), []byte(value))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess
	}

	err = instance.PutUint32(uint64(retPtr), uint32(addr))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess
	}

	err = instance.PutUint32(uint64(retSize), uint32(len(value)))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess
	}

	return api.WasmResultOk
}

func copyBytesIntoInstance(instance api.WasmInstance, value []byte, retPtr int32, retSize int32) api.WasmResult {
	addr, err := instance.Malloc(int32(len(value)))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess
	}

	err = instance.PutMemory(addr, uint64(len(value)), value)
	if err != nil {
		return api.WasmResultInvalidMemoryAccess
	}

	err = instance.PutUint32(uint64(retPtr), uint32(addr))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess
	}

	err = instance.PutUint32(uint64(retSize), uint32(len(value)))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess
	}

	return api.WasmResultOk
}

func getContextHandler(instance api.WasmInstance) api.ContextHandler {
	if v := instance.GetData(); v != nil {
		if im, ok := v.(api.ContextHandler); ok {
			return im
		}
	}

	return nil
}

func getImportHandler(instance api.WasmInstance) api.ImportsHandler {
	if ctx := getContextHandler(instance); ctx != nil {
		if im := ctx.GetImports(); im != nil {
			return im
		}
	}

	return &abi.DefaultImportsHandler{}
}
