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

package api

type ImportsHandler interface { //nolint:interfacebloat
	// Logging
	Log(level LogLevel, msg string) WasmResult
	GetLogLevel() LogLevel

	// System
	GetRootContextID() int32
	SetEffectiveContextID(contextID int32) WasmResult
	Done() WasmResult

	// Time
	SetTickPeriodMilliseconds(tickPeriodMilliseconds int32) WasmResult
	GetCurrentTimeNanoseconds() (int32, WasmResult)

	GetStatus() (int32, string, WasmResult)

	// Configuration
	GetVmConfig() IoBuffer
	GetPluginConfig() IoBuffer

	// Metric
	DefineMetric(metricType MetricType, name string) (int32, WasmResult)
	IncrementMetric(metricID int32, offset int64) WasmResult
	RecordMetric(metricID int32, value int64) WasmResult
	GetMetric(metricID int32) (int64, WasmResult)

	// State accessors
	GetProperty(key string) (string, WasmResult)
	SetProperty(key string, value string) WasmResult

	// L4
	GetDownStreamData() IoBuffer
	GetUpstreamData() IoBuffer
	ResumeDownstream() WasmResult
	ResumeUpstream() WasmResult
	CloseDownstream() WasmResult
	CloseUpstream() WasmResult

	// HTTP request
	GetHttpRequestHeader() HeaderMap
	GetHttpRequestBody() IoBuffer
	GetHttpRequestTrailer() HeaderMap
	ResumeHttpRequest() WasmResult
	CloseHttpRequest() WasmResult

	// HTTP response
	GetHttpResponseHeader() HeaderMap
	GetHttpResponseBody() IoBuffer
	GetHttpResponseTrailer() HeaderMap
	ResumeHttpResponse() WasmResult
	CloseHttpResponse() WasmResult
	SendHttpResp(respCode int32, respCodeDetail IoBuffer, respBody IoBuffer, additionalHeaderMap HeaderMap, grpcCode int32) WasmResult

	// HTTP call out
	HttpCall(url string, headers HeaderMap, body IoBuffer, trailer HeaderMap, timeoutMilliseconds int32) (int32, WasmResult)
	GetHttpCallResponseHeaders() HeaderMap
	GetHttpCallResponseBody() IoBuffer
	GetHttpCallResponseTrailer() HeaderMap

	// gRPC
	OpenGrpcStream(grpcService string, serviceName string, method string) (int32, WasmResult)
	SendGrpcCallMsg(token int32, data IoBuffer, endOfStream int32) WasmResult
	CancelGrpcCall(token int32) WasmResult
	CloseGrpcCall(token int32) WasmResult

	GrpcCall(grpcService string, serviceName string, method string, data IoBuffer, timeoutMilliseconds int32) (int32, WasmResult)
	GetGrpcReceiveInitialMetaData() HeaderMap
	GetGrpcReceiveBuffer() IoBuffer
	GetGrpcReceiveTrailerMetaData() HeaderMap

	// foreign
	CallForeignFunction(funcName string, param []byte) ([]byte, WasmResult)
	GetFuncCallData() IoBuffer

	// Shared data
	GetSharedData(key string) (string, uint32, WasmResult)
	SetSharedData(key string, value string, cas uint32) WasmResult

	// Shared queue
	RegisterSharedQueue(queueName string) (uint32, WasmResult)
	ResolveSharedQueue(queueName string) (uint32, WasmResult)
	EnqueueSharedQueue(queueID uint32, data string) WasmResult
	DequeueSharedQueue(queueID uint32) (string, WasmResult)

	// for golang host environment
	// Wait until async call return, eg. sync http call in golang
	Wait() Action
}
