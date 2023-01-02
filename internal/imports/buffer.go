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

package imports

import (
	"context"

	"mosn.io/proxy-wasm-go-host/api"
)

func GetBuffer(instance api.WasmInstance, bufferType api.BufferType) api.IoBuffer {
	im := getImportHandler(instance)

	switch bufferType {
	case api.BufferTypeHttpRequestBody:
		return im.GetHttpRequestBody()
	case api.BufferTypeHttpResponseBody:
		return im.GetHttpResponseBody()
	case api.BufferTypeDownstreamData:
		return im.GetDownStreamData()
	case api.BufferTypeUpstreamData:
		return im.GetUpstreamData()
	case api.BufferTypeHttpCallResponseBody:
		return im.GetHttpCallResponseBody()
	case api.BufferTypeGrpcReceiveBuffer:
		return im.GetGrpcReceiveBuffer()
	case api.BufferTypePluginConfiguration:
		return im.GetPluginConfig()
	case api.BufferTypeVmConfiguration:
		return im.GetVmConfig()
	case api.BufferTypeCallData:
		return im.GetFuncCallData()
	}

	return nil
}

func (h *host) ProxyGetBufferBytes(ctx context.Context, bufferType int32, start int32, length int32, returnDataPtr int32, returnDataSize int32) int32 {
	if api.BufferType(bufferType) > api.BufferTypeMax {
		return api.WasmResultBadArgument.Int32()
	}

	instance := h.Instance
	buf := GetBuffer(instance, api.BufferType(bufferType))
	if buf == nil {
		return api.WasmResultNotFound.Int32()
	}

	if start > start+length {
		return api.WasmResultBadArgument.Int32()
	}

	if start+length > int32(buf.Len()) {
		length = int32(buf.Len()) - start
	}

	addr, err := instance.Malloc(length)
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	err = instance.PutMemory(addr, uint64(length), buf.Bytes())
	if err != nil {
		return api.WasmResultInternalFailure.Int32()
	}

	err = instance.PutUint32(uint64(returnDataPtr), uint32(addr))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	err = instance.PutUint32(uint64(returnDataSize), uint32(length))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

func (h *host) ProxySetBufferBytes(ctx context.Context, bufferType int32, start int32, length int32, dataPtr int32, dataSize int32) int32 {
	if api.BufferType(bufferType) > api.BufferTypeMax {
		return api.WasmResultBadArgument.Int32()
	}

	instance := h.Instance
	buf := GetBuffer(instance, api.BufferType(bufferType))
	if buf == nil {
		return api.WasmResultNotFound.Int32()
	}

	content, err := instance.GetMemory(uint64(dataPtr), uint64(dataSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	switch {
	case start == 0:
		if length == 0 || int(length) >= buf.Len() {
			buf.Drain(buf.Len())
			_, err = buf.Write(content)
		} else {
			return api.WasmResultBadArgument.Int32()
		}
	case int(start) >= buf.Len():
		_, err = buf.Write(content)
	default:
		return api.WasmResultBadArgument.Int32()
	}

	if err != nil {
		return api.WasmResultInternalFailure.Int32()
	}

	return api.WasmResultOk.Int32()
}

func (h *host) ProxyGetBufferStatus(ctx context.Context, bufferType int32, lengthPtr int32, flagsPtr int32) int32 {
	if api.BufferType(bufferType) > api.BufferTypeMax {
		return api.WasmResultBadArgument.Int32()
	}

	instance := h.Instance
	buf := GetBuffer(instance, api.BufferType(bufferType))
	if buf == nil {
		return api.WasmResultNotFound.Int32()
	}

	if err := instance.PutUint32(uint64(lengthPtr), uint32(buf.Len())); err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	if err := instance.PutUint32(uint64(flagsPtr), 0); err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}
