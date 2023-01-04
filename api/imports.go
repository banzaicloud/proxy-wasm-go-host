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

package api

type Imports interface { //nolint:interfacebloat
	// Logging
	ProxyLog(logLevel int32, messagePtr int32, messageSize int32) int32
	ProxyGetLogLevel(logLevelPtr int32) int32

	// Timer (will be set for the root context, e.g. onStart, onTick).
	ProxySetTickPeriodMilliseconds(tickPeriodMillisecondsPtr int32) int32
	// Time
	ProxyGetCurrentTimeNanoseconds(resultUint64Ptr int32) int32

	// Results status details for any previous ABI call and onGrpcClose.
	ProxyGetStatus(statusCodePtr int32, returnStatusDetailPtr int32, returnStatusDetailSize int32) int32

	// System
	ProxySetEffectiveContext(contextID int32) int32
	ProxyDone() int32
	ProxyCallForeignFunction(funcNamePtr int32, funcNameSize int32, paramPtr int32, paramSize int32, returnDataPtr int32, returnSize int32) int32

	// State accessors
	ProxyGetProperty(keyPtr int32, keySize int32, returnValuePtr int32, returnValueSize int32) int32
	ProxySetProperty(keyPtr int32, keySize int32, valuePtr int32, valueSize int32) int32

	// Continue/Close/Reply
	ProxyContinueStream(streamType int32) int32
	ProxyCloseStream(streamType int32) int32
	ProxySendLocalResponse(statusCode int32, statusCodeDetailsPtr int32, statusCodeDetailsSize int32, bodyPtr int32, bodySize int32, headersPtr int32, headersSize int32, grpcStatus int32) int32

	// Headers/Trailers/Metadata Maps
	ProxyAddHeaderMapValue(mapType int32, keyPtr int32, keySize int32, valuePtr int32, valueSize int32) int32
	ProxyGetHeaderMapValue(mapType MapType, keyPtr int32, keySize int32, returnValuePtr int32, returnValueSize int32) int32
	ProxyGetHeaderMapPairs(mapType MapType, returnDataPtr int32, returnDataSize int32) int32
	ProxySetHeaderMapPairs(mapType MapType, dataPtr int32, dataSize int32) int32
	ProxyReplaceHeaderMapValue(mapType int32, keyPtr int32, keySize int32, valuePtr int32, valueSize int32) int32
	ProxyRemoveHeaderMapValue(mapType int32, keyPtr int32, keySize int32) int32
	ProxyGetHeaderMapSize(mapType int32, sizePtr int32) int32

	// Shared data
	// Returns: Ok, NotFound
	ProxyGetSharedData(keyPtr int32, keySize int32, returnValuePtr int32, returnValueSize int32, returnCasPtr int32) int32
	// If cas != 0 and cas != the current cas for 'key' return false, otherwise set
	// the value and return true.
	// Returns: Ok, CasMismatch
	ProxySetSharedData(keyPtr int32, keySize int32, valuePtr int32, valueSize int32, cas int32) int32

	// Shared queue
	// Note: Registering the same queue_name will overwrite the old registration
	// while preseving any pending data. Consequently it should typically be
	// followed by a call to proxy_dequeue_shared_queue. Returns: Ok
	ProxyRegisterSharedQueue(queueNamePtr int32, queueNameSize int32, tokenIDPtr int32) int32
	// Returns: Ok, NotFound
	ProxyResolveSharedQueue(vmIDPtr int32, vmIDSize int32, queueNamePtr int32, queueNameSize int32, tokenIDPtr int32) int32
	// Returns false if the queue was not found and the data was not enqueued.
	ProxyEnqueueSharedQueue(tokenID int32, dataPtr int32, dataSize int32) int32
	// Returns Ok, Empty, NotFound (token not registered).
	ProxyDequeueSharedQueue(tokenID int32, returnValuePtr int32, returnValueSize int32) int32

	// Buffer
	ProxyGetBufferBytes(bufferType int32, start int32, length int32, returnDataPtr int32, returnDataSize int32) int32
	ProxyGetBufferStatus(bufferType int32, lengthPtr int32, flagsPtr int32) int32
	ProxySetBufferBytes(bufferType int32, start int32, length int32, dataPtr int32, dataSize int32) int32

	// Metrics
	ProxyDefineMetric(metricType int32, namePtr int32, nameSize int32, returnMetricId int32) int32
	ProxyIncrementMetric(metricId int32, offset int64) int32
	ProxyRecordMetric(metricId int32, value int64) int32
	ProxyGetMetric(metricId int32, resultUint64Ptr int32) int32

	// HTTP
	ProxyHttpCall(uriPtr int32, uriSize int32, headerPairsPtr int32, headerPairsSize int32, bodyPtr int32, bodySize int32, trailerPairsPtr int32, trailerPairsSize int32, timeoutMilliseconds int32, calloutIDPtr int32) int32

	// gRPC
	ProxyGrpcCall(grpcServiceData int32, grpcServiceSize int32, serviceNameData int32, serviceNameSize int32, methodName int32, methodSize int32, initialMetadataPtr int32, initialMetadataSize int32, grpcMessage int32, grpcMessageSize int32, timeoutMilliseconds int32, returnCalloutID int32) int32
	ProxyGrpcStream(grpcServiceData int32, grpcServiceSize int32, serviceNameData int32, serviceNameSize int32, methodName int32, methodSize int32, initialMetadataPtr int32, initialMetadataSize int32, grpcMessage int32, grpcMessageSize int32, returnStreamID int32) int32
	ProxyGrpcSend(streamID int32, messagePtr int32, messageSize int32, endOfStream int32) int32
	ProxyGrpcCancel(calloutID int32) int32
	ProxyGrpcClose(calloutID int32) int32
}
