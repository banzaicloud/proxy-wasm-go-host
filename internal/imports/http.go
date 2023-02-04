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
	"bytes"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
	"github.com/banzaicloud/proxy-wasm-go-host/pkg/utils"
)

func (h *host) ProxyHttpCall(
	uriPtr int32,
	uriSize int32,
	headerPairsPtr int32,
	headerPairsSize int32,
	bodyPtr int32,
	bodySize int32,
	trailerPairsPtr int32,
	trailerPairsSize int32,
	timeoutMilliseconds int32,
	calloutIDPtr int32,
) int32 {
	instance := h.Instance
	url, err := instance.GetMemory(uint64(uriPtr), uint64(uriSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	headerMapData, err := instance.GetMemory(uint64(headerPairsPtr), uint64(headerPairsSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	headerMap := utils.DecodeMap(headerMapData)

	body, err := instance.GetMemory(uint64(bodyPtr), uint64(bodySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	trailerMapData, err := instance.GetMemory(uint64(trailerPairsPtr), uint64(trailerPairsSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	trailerMap := utils.DecodeMap(trailerMapData)

	ih := getImportHandler(instance)

	calloutID, res := ih.HttpCall(
		string(url),
		utils.CommonHeader(headerMap),
		bytes.NewBuffer(body),
		utils.CommonHeader(trailerMap),
		timeoutMilliseconds,
	)
	if res != api.WasmResultOk {
		return res.Int32()
	}

	err = instance.PutUint32(uint64(calloutIDPtr), uint32(calloutID))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}
