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

var _ api.Imports = &host{}

type host struct {
	Instance api.WasmInstance
}

func HostFunctions(instance api.WasmInstance) map[string]interface{} {
	h := &host{Instance: instance}

	return map[string]interface{}{
		"proxy_log":                          h.ProxyLog,
		"proxy_get_log_level":                h.ProxyGetLogLevel,
		"proxy_set_tick_period_milliseconds": h.ProxySetTickPeriodMilliseconds,
		"proxy_get_current_time_nanoseconds": h.ProxyGetCurrentTimeNanoseconds,
		"proxy_get_status":                   h.ProxyGetStatus,
		"proxy_get_property":                 h.ProxyGetProperty,
		"proxy_set_property":                 h.ProxySetProperty,
		"proxy_continue_stream":              h.ProxyContinueStream,
		"proxy_close_stream":                 h.ProxyCloseStream,
		"proxy_add_header_map_value":         h.ProxyAddHeaderMapValue,
		"proxy_get_header_map_value":         h.ProxyGetHeaderMapValue,
		"proxy_get_header_map_pairs":         h.ProxyGetHeaderMapPairs,
		"proxy_set_header_map_pairs":         h.ProxySetHeaderMapPairs,
		"proxy_get_header_map_size":          h.ProxyGetHeaderMapSize,
		"proxy_get_shared_data":              h.ProxyGetSharedData,
		"proxy_set_shared_data":              h.ProxySetSharedData,
		"proxy_register_shared_queue":        h.ProxyRegisterSharedQueue,
		"proxy_resolve_shared_queue":         h.ProxyResolveSharedQueue,
		"proxy_enqueue_shared_queue":         h.ProxyEnqueueSharedQueue,
		"proxy_dequeue_shared_queue":         h.ProxyDequeueSharedQueue,
		"proxy_replace_header_map_value":     h.ProxyReplaceHeaderMapValue,
		"proxy_remove_header_map_value":      h.ProxyRemoveHeaderMapValue,
		"proxy_get_buffer_bytes":             h.ProxyGetBufferBytes,
		"proxy_get_buffer_status":            h.ProxyGetBufferStatus,
		"proxy_set_buffer_bytes":             h.ProxySetBufferBytes,
		"proxy_http_call":                    h.ProxyHttpCall,
		"proxy_define_metric":                h.ProxyDefineMetric,
		"proxy_increment_metric":             h.ProxyIncrementMetric,
		"proxy_record_metric":                h.ProxyRecordMetric,
		"proxy_get_metric":                   h.ProxyGetMetric,
		"proxy_grpc_call":                    h.ProxyGrpcCall,
		"proxy_grpc_stream":                  h.ProxyGrpcStream,
		"proxy_grpc_send":                    h.ProxyGrpcSend,
		"proxy_grpc_close":                   h.ProxyGrpcClose,
		"proxy_grpc_cancel":                  h.ProxyGrpcCancel,
		"proxy_set_effective_context":        h.ProxySetEffectiveContext,
		"proxy_done":                         h.ProxyDone,
		"proxy_call_foreign_function":        h.ProxyCallForeignFunction,
		"proxy_send_local_response":          h.ProxySendLocalResponse,
	}
}

// Logging

func (h *host) ProxyLog(logLevel int32, messagePtr int32, messageSize int32) int32 {
	instance := h.Instance
	ih := getImportHandler(h.Instance)

	logContent, err := instance.GetMemory(uint64(messagePtr), uint64(messageSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return ih.Log(api.LogLevel(logLevel), string(logContent)).Int32()
}

func (h *host) ProxyGetLogLevel(logLevelPtr int32) int32 {
	instance := h.Instance
	ih := getImportHandler(h.Instance)

	if err := instance.PutUint32(uint64(logLevelPtr), uint32(ih.GetLogLevel())); err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

// Results status details for any previous ABI call and onGrpcClose.
func (h *host) ProxyGetStatus(statusCodePtr int32, returnStatusDetailPtr int32, returnStatusDetailSize int32) int32 {
	instance := h.Instance
	ih := getImportHandler(h.Instance)

	statusCode, statusMessage, res := ih.GetStatus()
	if res != api.WasmResultOk {
		return res.Int32()
	}

	if err := instance.PutUint32(uint64(statusCodePtr), uint32(statusCode)); err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return copyIntoInstance(instance, statusMessage, returnStatusDetailPtr, returnStatusDetailSize).Int32()
}

// Timer (will be set for the root context, e.g. onStart, onTick).

func (h *host) ProxySetTickPeriodMilliseconds(tickPeriodMilliseconds int32) int32 {
	ih := getImportHandler(h.Instance)

	return ih.SetTickPeriodMilliseconds(tickPeriodMilliseconds).Int32()
}

// Time

func (h *host) ProxyGetCurrentTimeNanoseconds(resultUint64Ptr int32) int32 {
	instance := h.Instance
	ih := getImportHandler(instance)

	nano, res := ih.GetCurrentTimeNanoseconds()
	if res != api.WasmResultOk {
		return res.Int32()
	}

	if err := instance.PutUint32(uint64(resultUint64Ptr), uint32(nano)); err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

// System

func (h *host) ProxySetEffectiveContext(contextID int32) int32 {
	ih := getImportHandler(h.Instance)

	return ih.SetEffectiveContextID(contextID).Int32()
}

func (h *host) ProxyDone() int32 {
	ih := getImportHandler(h.Instance)

	return ih.Done().Int32()
}

func (h *host) ProxyCallForeignFunction(funcNamePtr int32, funcNameSize int32,
	paramPtr int32, paramSize int32, returnData int32, returnSize int32,
) int32 {
	instance := h.Instance
	funcName, err := instance.GetMemory(uint64(funcNamePtr), uint64(funcNameSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	param, err := instance.GetMemory(uint64(paramPtr), uint64(paramSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	ih := getImportHandler(instance)

	ret, res := ih.CallForeignFunction(string(funcName), param)
	if res != api.WasmResultOk {
		return res.Int32()
	}

	return copyBytesIntoInstance(instance, ret, returnData, returnSize).Int32()
}

// State accessors

func (h *host) ProxyGetProperty(keyPtr int32, keySize int32, returnValuePtr int32, returnValueSize int32) int32 {
	instance := h.Instance
	key, err := instance.GetMemory(uint64(keyPtr), uint64(keySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	ih := getImportHandler(instance)

	value, res := ih.GetProperty(string(key))
	if res != api.WasmResultOk {
		return res.Int32()
	}

	return copyIntoInstance(instance, value, returnValuePtr, returnValueSize).Int32()
}

func (h *host) ProxySetProperty(keyPtr int32, keySize int32, valuePtr int32, valueSize int32) int32 {
	instance := h.Instance
	key, err := instance.GetMemory(uint64(keyPtr), uint64(keySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	value, err := instance.GetMemory(uint64(valuePtr), uint64(valueSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	ih := getImportHandler(instance)

	return ih.SetProperty(string(key), string(value)).Int32()
}

// Continue/Close/Reply

func (h *host) ProxyContinueStream(streamType int32) int32 {
	ih := getImportHandler(h.Instance)

	switch api.StreamType(streamType) {
	case api.StreamTypeDownstream:
		return ih.ResumeDownstream().Int32()
	case api.StreamTypeUpstream:
		return ih.ResumeUpstream().Int32()
	case api.StreamTypeHttpRequest:
		return ih.ResumeHttpRequest().Int32()
	case api.StreamTypeHttpResponse:
		return ih.ResumeHttpResponse().Int32()
	}

	return api.WasmResultBadArgument.Int32()
}

func (h *host) ProxyCloseStream(streamType int32) int32 {
	ih := getImportHandler(h.Instance)

	switch api.StreamType(streamType) {
	case api.StreamTypeDownstream:
		return ih.CloseDownstream().Int32()
	case api.StreamTypeUpstream:
		return ih.CloseUpstream().Int32()
	case api.StreamTypeHttpRequest:
		return ih.CloseHttpRequest().Int32()
	case api.StreamTypeHttpResponse:
		return ih.CloseHttpResponse().Int32()
	}

	return api.WasmResultBadArgument.Int32()
}

func (h *host) ProxySendLocalResponse(
	statusCode int32,
	statusCodeDetailsPtr int32,
	statusCodeDetailsSize int32,
	bodyPtr int32,
	bodySize int32,
	headersPtr int32,
	headersSize int32,
	grpcStatus int32,
) int32 {
	instance := h.Instance
	statusCodeDetail, err := instance.GetMemory(uint64(statusCodeDetailsPtr), uint64(statusCodeDetailsSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	respBody, err := instance.GetMemory(uint64(bodyPtr), uint64(bodySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	headers, err := instance.GetMemory(uint64(headersPtr), uint64(headersSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	additionalHeaderMap := utils.DecodeMap(headers)

	ih := getImportHandler(instance)

	return ih.SendHttpResp(statusCode,
		bytes.NewBuffer(statusCodeDetail),
		bytes.NewBuffer(respBody),
		utils.CommonHeader(additionalHeaderMap), grpcStatus).Int32()
}
