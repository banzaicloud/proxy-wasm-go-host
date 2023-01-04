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
	"github.com/banzaicloud/proxy-wasm-go-host/api"
	"github.com/banzaicloud/proxy-wasm-go-host/pkg/utils"
)

// u32 is the fixed size of a uint32 in little-endian encoding.
const u32Len = 4

func GetMap(instance api.WasmInstance, mapType api.MapType) api.HeaderMap {
	ih := getImportHandler(instance)

	switch mapType {
	case api.MapTypeHttpRequestHeaders:
		return ih.GetHttpRequestHeader()
	case api.MapTypeHttpRequestTrailers:
		return ih.GetHttpRequestTrailer()
	case api.MapTypeHttpResponseHeaders:
		return ih.GetHttpResponseHeader()
	case api.MapTypeHttpResponseTrailers:
		return ih.GetHttpResponseTrailer()
	case api.MapTypeGrpcReceiveInitialMetadata:
		return ih.GetGrpcReceiveInitialMetaData()
	case api.MapTypeGrpcReceiveTrailingMetadata:
		return ih.GetGrpcReceiveTrailerMetaData()
	case api.MapTypeHttpCallResponseHeaders:
		return ih.GetHttpCallResponseHeaders()
	case api.MapTypeHttpCallResponseTrailers:
		return ih.GetHttpCallResponseTrailer()
	}

	return nil
}

// Headers/Trailers/Metadata Maps

func (h *host) ProxyAddHeaderMapValue(mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	instance := h.Instance
	headerMap := GetMap(instance, mapType)
	if headerMap == nil {
		return api.WasmResultNotFound.Int32()
	}

	key, err := instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	value, err := instance.GetMemory(uint64(valueDataPtr), uint64(valueSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	headerMap.Set(string(key), string(value))

	return api.WasmResultOk.Int32()
}

func (h *host) ProxyGetHeaderMapValue(mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	instance := h.Instance
	headerMap := GetMap(instance, mapType)
	if headerMap == nil {
		return api.WasmResultNotFound.Int32()
	}

	key, err := instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	value, ok := headerMap.Get(string(key))
	if !ok {
		return api.WasmResultNotFound.Int32()
	}

	return copyIntoInstance(instance, value, valueDataPtr, valueSize).Int32()
}

func (h *host) ProxyGetHeaderMapPairs(mapType int32, returnDataPtr int32, returnDataSize int32) int32 {
	instance := h.Instance
	header := GetMap(instance, mapType)
	if header == nil {
		return api.WasmResultNotFound.Int32()
	}

	cloneMap := make(map[string]string)
	totalBytesLen := u32Len
	header.Range(func(key, value string) bool {
		cloneMap[key] = value
		totalBytesLen += u32Len + u32Len               // keyLen + valueLen
		totalBytesLen += len(key) + 1 + len(value) + 1 // key + \0 + value + \0

		return true
	})

	addr, err := instance.Malloc(int32(totalBytesLen))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	err = instance.PutUint32(addr, uint32(len(cloneMap)))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	lenPtr := addr + u32Len
	dataPtr := lenPtr + uint64((u32Len+u32Len)*len(cloneMap))

	for k, v := range cloneMap {
		_ = instance.PutUint32(lenPtr, uint32(len(k)))
		lenPtr += u32Len
		_ = instance.PutUint32(lenPtr, uint32(len(v)))
		lenPtr += u32Len

		_ = instance.PutMemory(dataPtr, uint64(len(k)), []byte(k))
		dataPtr += uint64(len(k))
		_ = instance.PutByte(dataPtr, 0)
		dataPtr++

		_ = instance.PutMemory(dataPtr, uint64(len(v)), []byte(v))
		dataPtr += uint64(len(v))
		_ = instance.PutByte(dataPtr, 0)
		dataPtr++
	}

	err = instance.PutUint32(uint64(returnDataPtr), uint32(addr))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	err = instance.PutUint32(uint64(returnDataSize), uint32(totalBytesLen))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

func (h *host) ProxySetHeaderMapPairs(mapType int32, ptr int32, size int32) int32 {
	instance := h.Instance
	headerMap := GetMap(instance, mapType)
	if headerMap == nil {
		return api.WasmResultNotFound.Int32()
	}

	newMapContent, err := instance.GetMemory(uint64(ptr), uint64(size))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	newMap := utils.DecodeMap(newMapContent)

	for k, v := range newMap {
		headerMap.Set(k, v)
	}

	return api.WasmResultOk.Int32()
}

func (h *host) ProxyReplaceHeaderMapValue(mapType int32, keyDataPtr int32, keySize int32, valueDataPtr int32, valueSize int32) int32 {
	instance := h.Instance
	headerMap := GetMap(instance, mapType)
	if headerMap == nil {
		return api.WasmResultNotFound.Int32()
	}

	key, err := instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	value, err := instance.GetMemory(uint64(valueDataPtr), uint64(valueSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(value) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	headerMap.Set(string(key), string(value))

	return api.WasmResultOk.Int32()
}

func (h *host) ProxyRemoveHeaderMapValue(mapType int32, keyDataPtr int32, keySize int32) int32 {
	instance := h.Instance
	headerMap := GetMap(instance, mapType)
	if headerMap == nil {
		return api.WasmResultNotFound.Int32()
	}

	key, err := instance.GetMemory(uint64(keyDataPtr), uint64(keySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	headerMap.Del(string(key))

	return api.WasmResultOk.Int32()
}

func (h *host) ProxyGetHeaderMapSize(mapType int32, sizePtr int32) int32 {
	instance := h.Instance
	headerMap := GetMap(instance, mapType)
	if headerMap == nil {
		return api.WasmResultNotFound.Int32()
	}

	if err := instance.PutUint32(uint64(sizePtr), uint32(headerMap.ByteSize())); err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}
