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

type KVStore interface {
	HeaderMap
	SetCAS(key, value string, cas bool) bool
	DelCAS(key string, cas bool) bool
}

type LogLevel int32

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelCritical
)

type Result int32

const (
	ResultOk Result = iota
	ResultEmpty
	ResultNotFound
	ResultNotAllowed
	ResultBadArgument
	ResultInvalidMemoryAccess
	ResultInvalidOperation
	ResultCompareAndSwapMismatch
	ResultUnimplemented Result = 12
)

type WasmResult int32

const (
	WasmResultOk                   WasmResult = iota
	WasmResultNotFound                        // The result could not be found, e.g. a provided key did not appear in a table.
	WasmResultBadArgument                     // An argument was bad, e.g. did not conform to the required range.
	WasmResultSerializationFailure            // A protobuf could not be serialized.
	WasmResultParseFailure                    // A protobuf could not be parsed.
	WasmResultBadExpression                   // A provided expression (e.g. "foo.bar") was illegal or unrecognized.
	WasmResultInvalidMemoryAccess             // A provided memory range was not legal.
	WasmResultEmpty                           // Data was requested from an empty container.
	WasmResultCasMismatch                     // The provided CAS did not match that of the stored data.
	WasmResultResultMismatch                  // Returned result was unexpected, e.g. of the incorrect size.
	WasmResultInternalFailure                 // Internal failure: trying check logs of the surrounding system.
	WasmResultBrokenConnection                // The connection/stream/pipe was broken/closed unexpectedly.
	WasmResultUnimplemented                   // Feature not implemented.
)

func (wasmResult WasmResult) Int32() int32 {
	return int32(wasmResult)
}

type Action int32

const (
	ActionContinue Action = iota
	ActionPause
)

type StreamType int32

const (
	StreamTypeHttpRequest StreamType = iota
	StreamTypeHttpResponse
	StreamTypeDownstream
	StreamTypeUpstream
)

type ContextType int32

const (
	ContextTypeHttpContext ContextType = iota
	ContextTypeStreamContext
)

type BufferType int32

const (
	BufferTypeHttpRequestBody BufferType = iota
	BufferTypeHttpResponseBody
	BufferTypeDownstreamData
	BufferTypeUpstreamData
	BufferTypeHttpCallResponseBody
	BufferTypeGrpcReceiveBuffer
	BufferTypeVmConfiguration
	BufferTypePluginConfiguration
	BufferTypeCallData
	BufferTypeMax BufferType = 8
)

type MapType = int32

const (
	MapTypeHttpRequestHeaders MapType = iota
	MapTypeHttpRequestTrailers
	MapTypeHttpResponseHeaders
	MapTypeHttpResponseTrailers
	MapTypeGrpcReceiveInitialMetadata
	MapTypeGrpcReceiveTrailingMetadata
	MapTypeHttpCallResponseHeaders
	MapTypeHttpCallResponseTrailers
	MapTypeMax MapType = 7
)

type PeerType int32

const (
	PeerTypeUnknown PeerType = iota
	PeerTypeLocal            // Close initiated by the proxy.
	PeerTypeRemote           // Close initiated by the peer.
)

type MetricType int32

const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
	MetricTypeHistogram
	MetricTypeMax MetricType = 2
)
