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

package abi

import (
	"errors"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

const ProxyWasmABI_0_1_0 string = "proxy_abi_version_0_1_0"
const ProxyWasmABI_0_2_1 string = "proxy_abi_version_0_2_1"

func NewContext(importsHandler api.ImportsHandler, instance api.WasmInstance) api.ABIContext {
	return &context{
		imports:  importsHandler,
		instance: instance,
	}
}

type context struct {
	imports  api.ImportsHandler
	instance api.WasmInstance
}

func (a *context) Name() string {
	return ProxyWasmABI_0_2_1
}

func (a *context) GetExports() api.Exports {
	return a
}

func (a *context) GetImports() api.ImportsHandler {
	return a.imports
}

func (a *context) SetImports(imports api.ImportsHandler) {
	a.imports = imports
}

func (a *context) GetInstance() api.WasmInstance {
	return a.instance
}

func (a *context) SetInstance(instance api.WasmInstance) {
	a.instance = instance
}

//

var ErrInvalidResult = errors.New("invalid result")

func (a *context) CallWasmFunction(funcName string, args ...interface{}) (interface{}, api.Action, error) {
	ff, err := a.instance.GetExportsFunc(funcName)
	if err != nil {
		return nil, api.ActionContinue, err
	}

	res, err := ff.Call(args...)
	if err != nil {
		a.instance.HandleError(err)

		return nil, api.ActionContinue, err
	}

	// if we have sync call, e.g. HttpCall, then unlock the wasm instance and wait until it resp
	Action := a.imports.Wait()

	return res, Action, nil
}

// Configuration

func (a *context) ProxyOnVmStart(rootContextID int32, vmConfigurationSize int32) (bool, error) {
	res, _, err := a.CallWasmFunction("proxy_on_vm_start", rootContextID, vmConfigurationSize)
	if err != nil {
		return false, err
	}

	if v, ok := res.(int32); ok {
		return v == 1, nil
	}

	return false, ErrInvalidResult
}

func (a *context) ProxyOnConfigure(rootContextID int32, configurationSize int32) (bool, error) {
	res, _, err := a.CallWasmFunction("proxy_on_configure", rootContextID, configurationSize)
	if err != nil {
		return false, err
	}

	if v, ok := res.(int32); ok {
		return v == 1, nil
	}

	return false, ErrInvalidResult
}

// Misc

func (a *context) ProxyOnLog(contextID int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_log", contextID)

	return err
}

func (a *context) ProxyOnTick(rootContextID int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_tick", rootContextID)

	return err
}

func (a *context) ProxyOnHttpCallResponse(contextID int32, tokenID int32, headerCount int32, bodySize int32, trailerCount int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_http_call_response", contextID, tokenID, headerCount, bodySize, trailerCount)

	return err
}

func (a *context) ProxyOnQueueReady(rootContextID int32, queueID int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_queue_ready", rootContextID, queueID)

	return err
}

// Context

func (a *context) ProxyOnContextCreate(contextID int32, rootContextID int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_context_create", contextID, rootContextID)

	return err
}

func (a *context) ProxyOnDone(contextID int32) (bool, error) {
	res, _, err := a.CallWasmFunction("proxy_on_done", contextID)
	if err != nil {
		return false, err
	}

	if v, ok := res.(int32); ok {
		return v == 1, nil
	}

	return false, ErrInvalidResult
}

func (a *context) ProxyOnDelete(contextID int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_delete", contextID)

	return err
}

// L4

func (a *context) ProxyOnNewConnection(contextID int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_new_connection", contextID)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}

func (a *context) ProxyOnDownstreamData(contextID int32, dataSize int32, endOfStream int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_downstream_data", contextID, dataSize, endOfStream)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}

func (a *context) ProxyOnDownstreamConnectionClose(contextID int32, peerType api.PeerType) error {
	_, _, err := a.CallWasmFunction("proxy_on_downstream_connection_close", contextID, peerType)

	return err
}

func (a *context) ProxyOnUpstreamData(contextID int32, dataSize int32, endOfStream int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_upstream_data", contextID, dataSize, endOfStream)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}

func (a *context) ProxyOnUpstreamConnectionClose(contextID int32, peerType api.PeerType) error {
	_, _, err := a.CallWasmFunction("proxy_on_upstream_connection_close", contextID, peerType)

	return err
}

// gRPC

func (a *context) ProxyOnGrpcClose(contextID int32, calloutID int32, statusCode int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_grpc_close", contextID, calloutID, statusCode)

	return err
}

func (a *context) ProxyOnGrpcReceiveInitialMetadata(contextID int32, tokenID int32, headerCount int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_grpc_receive_initial_metadata", contextID, tokenID, headerCount)

	return err
}

func (a *context) ProxyOnGrpcReceiveTrailingMetadata(contextID int32, tokenID int32, trailerCount int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_grpc_receive_trailing_metadata", contextID, tokenID, trailerCount)

	return err
}

func (a *context) ProxyOnGrpcReceive(contextID int32, tokenID int32, responseSize int32) error {
	_, _, err := a.CallWasmFunction("proxy_on_grpc_receive", contextID, tokenID, responseSize)

	return err
}

// HTTP request

func (a *context) ProxyOnRequestBody(contextID int32, bodySize int32, endOfStream int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_request_body", contextID, bodySize, endOfStream)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}

func (a *context) ProxyOnRequestHeaders(contextID int32, headerCount int32, endOfStream int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_request_headers", contextID, headerCount, endOfStream)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}

func (a *context) ProxyOnRequestTrailers(contextID int32, trailerCount int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_request_trailers", contextID, trailerCount)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}

// HTTP response

func (a *context) ProxyOnResponseBody(contextID int32, bodySize int32, endOfStream int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_response_body", contextID, bodySize, endOfStream)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}

func (a *context) ProxyOnResponseHeaders(contextID int32, headerCount int32, endOfStream int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_response_headers", contextID, headerCount, endOfStream)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}

func (a *context) ProxyOnResponseTrailers(contextID int32, trailerCount int32) (api.Action, error) {
	_, Action, err := a.CallWasmFunction("proxy_on_response_trailers", contextID, trailerCount)
	if err != nil {
		return api.ActionPause, err
	}

	return Action, nil
}
