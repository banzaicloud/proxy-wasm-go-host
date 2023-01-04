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

type Exports interface { //nolint:interfacebloat
	// Configuration
	ProxyOnVmStart(rootContextID int32, vmConfigurationSize int32) (bool, error)
	ProxyOnConfigure(rootContextID int32, pluginConfigurationSize int32) (bool, error)

	// Misc
	ProxyOnLog(contextID int32) error
	ProxyOnTick(rootContextID int32) error
	ProxyOnHttpCallResponse(contextID int32, tokenID int32, headerCount int32, bodySize int32, trailerCount int32) error
	ProxyOnQueueReady(rootContextID int32, queueID int32) error

	// Context
	ProxyOnContextCreate(contextID int32, rootContextID int32) error
	ProxyOnDone(contextID int32) (bool, error)
	ProxyOnDelete(contextID int32) error

	// L4
	ProxyOnNewConnection(contextID int32) (Action, error)
	ProxyOnDownstreamData(contextID int32, dataSize int32, endOfStream int32) (Action, error)
	ProxyOnDownstreamConnectionClose(contextID int32, peerType int32) error
	ProxyOnUpstreamData(contextID int32, dataSize int32, endOfStream int32) (Action, error)
	ProxyOnUpstreamConnectionClose(contextID int32, peerType int32) error

	// gRPC
	ProxyOnGrpcClose(contextID int32, tokenID int32, statusCode int32) error
	ProxyOnGrpcReceiveInitialMetadata(contextID int32, tokenID int32, headerCount int32) error
	ProxyOnGrpcReceiveTrailingMetadata(contextID int32, tokenID int32, trailerCount int32) error
	ProxyOnGrpcReceive(contextID int32, tokenID int32, responseSize int32) error

	// HTTP request
	ProxyOnRequestBody(contextID int32, bodySize int32, endOfStream int32) (Action, error)
	ProxyOnRequestHeaders(contextID int32, headerCount int32, endOfStream int32) (Action, error)
	ProxyOnRequestTrailers(contextID int32, trailerCount int32) (Action, error)

	// HTTP response
	ProxyOnResponseBody(contextID int32, bodySize int32, endOfStream int32) (Action, error)
	ProxyOnResponseHeaders(contextID int32, headerCount int32, endOfStream int32) (Action, error)
	ProxyOnResponseTrailers(contextID int32, trailerCount int32) (Action, error)
}
