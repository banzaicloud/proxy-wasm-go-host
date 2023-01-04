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
	utils "github.com/banzaicloud/proxy-wasm-go-host/pkg/utils"
)

// gRPC

// TODO(@wayne): implement metadata
func (h *host) ProxyGrpcCall(
	grpcServiceData int32,
	grpcServiceSize int32,
	serviceNameData int32,
	serviceNameSize int32,
	methodName int32,
	methodSize int32,
	initialMetadataPtr int32,
	initialMetadataSize int32,
	grpcMessage int32,
	grpcMessageSize int32,
	timeoutMilliseconds int32,
	returnCalloutID int32,
) int32 {
	instance := h.Instance
	grpcService, err := instance.GetMemory(uint64(grpcServiceData), uint64(grpcServiceSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	serviceName, err := instance.GetMemory(uint64(serviceNameData), uint64(serviceNameSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	method, err := instance.GetMemory(uint64(methodName), uint64(methodSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	msg, err := instance.GetMemory(uint64(grpcMessage), uint64(grpcMessageSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	ih := getImportHandler(instance)

	calloutID, res := ih.GrpcCall(string(grpcService), string(serviceName), string(method),
		utils.NewIoBufferBytes(msg), timeoutMilliseconds)
	if res != api.WasmResultOk {
		return res.Int32()
	}

	err = instance.PutUint32(uint64(returnCalloutID), uint32(calloutID))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

// TODO(@wayne): implement metadata, message
func (h *host) ProxyGrpcStream(
	grpcServiceData int32,
	grpcServiceSize int32,
	serviceNameData int32,
	serviceNameSize int32,
	methodName int32,
	methodSize int32,
	initialMetadataPtr int32,
	initialMetadataSize int32,
	returnStreamID int32,
) int32 {
	instance := h.Instance
	grpcService, err := instance.GetMemory(uint64(grpcServiceData), uint64(grpcServiceSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	serviceName, err := instance.GetMemory(uint64(serviceNameData), uint64(serviceNameSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	method, err := instance.GetMemory(uint64(methodName), uint64(methodSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	ih := getImportHandler(instance)

	calloutID, res := ih.OpenGrpcStream(string(grpcService), string(serviceName), string(method))
	if res != api.WasmResultOk {
		return res.Int32()
	}

	err = instance.PutUint32(uint64(returnStreamID), uint32(calloutID))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

func (h *host) ProxyGrpcSend(streamID int32, messagePtr int32, messageSize int32, endOfStream int32) int32 {
	instance := h.Instance
	msg, err := instance.GetMemory(uint64(messagePtr), uint64(messageSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	ih := getImportHandler(instance)

	return ih.SendGrpcCallMsg(streamID, utils.NewIoBufferBytes(msg), endOfStream).Int32()
}

func (h *host) ProxyGrpcCancel(calloutID int32) int32 {
	instance := h.Instance
	ih := getImportHandler(instance)

	return ih.CancelGrpcCall(calloutID).Int32()
}

func (h *host) ProxyGrpcClose(calloutID int32) int32 {
	instance := h.Instance
	ih := getImportHandler(instance)

	return ih.CloseGrpcCall(calloutID).Int32()
}
