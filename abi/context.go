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

package abi

import (
	"errors"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

const ProxyWasmABI_0_1_0 string = "proxy_abi_version_0_1_0"
const ProxyWasmABI_0_2_0 string = "proxy_abi_version_0_2_0"
const ProxyWasmABI_0_2_1 string = "proxy_abi_version_0_2_1"

func NewContext(importsHandler api.ImportsHandler, instance api.WasmInstance) (api.ABIContext, error) {
	proxyOnDownstreamData, err := instance.GetExportsFunc("proxy_on_downstream_data")
	if err != nil {
		return nil, err
	}

	proxyOnUpstreamData, err := instance.GetExportsFunc("proxy_on_upstream_data")
	if err != nil {
		return nil, err
	}

	proxyOnUpstreamConnectionClose, err := instance.GetExportsFunc("proxy_on_upstream_connection_close")
	if err != nil {
		return nil, err
	}

	proxyOnDownstreamConnectionClose, err := instance.GetExportsFunc("proxy_on_downstream_connection_close")
	if err != nil {
		return nil, err
	}

	proxyOnLog, err := instance.GetExportsFunc("proxy_on_log")
	if err != nil {
		return nil, err
	}

	proxyOnTick, err := instance.GetExportsFunc("proxy_on_tick")
	if err != nil {
		return nil, err
	}

	proxyOnNewConnection, err := instance.GetExportsFunc("proxy_on_new_connection")
	if err != nil {
		return nil, err
	}

	proxyOnContextCreate, err := instance.GetExportsFunc("proxy_on_context_create")
	if err != nil {
		return nil, err
	}

	proxyOnDone, err := instance.GetExportsFunc("proxy_on_done")
	if err != nil {
		return nil, err
	}

	proxyOnDelete, err := instance.GetExportsFunc("proxy_on_delete")
	if err != nil {
		return nil, err
	}

	proxyOnQueueReady, err := instance.GetExportsFunc("proxy_on_queue_ready")
	if err != nil {
		return nil, err
	}

	proxyOnGrpcClose, _ := instance.GetExportsFunc("proxy_on_grpc_close")

	proxyOnGrpcReceiveInitialMetadata, _ := instance.GetExportsFunc("proxy_on_grpc_receive_initial_metadata")

	proxyOnGrpcReceiveTrailingMetadata, _ := instance.GetExportsFunc("proxy_on_grpc_receive_trailing_metadata")

	proxyOnGrpcReceive, _ := instance.GetExportsFunc("proxy_on_grpc_receive")

	proxyOnRequestBody, _ := instance.GetExportsFunc("proxy_on_request_body")

	proxyOnRequestHeaders, _ := instance.GetExportsFunc("proxy_on_request_headers")

	proxyOnRequestTrailers, _ := instance.GetExportsFunc("proxy_on_request_trailers")

	proxyOnResponseBody, _ := instance.GetExportsFunc("proxy_on_response_body")

	proxyOnResponseHeaders, _ := instance.GetExportsFunc("proxy_on_response_headers")

	proxyOnResponseTrailers, _ := instance.GetExportsFunc("proxy_on_response_trailers")

	proxyOnHttpCallResponse, _ := instance.GetExportsFunc("proxy_on_http_call_response")

	return &context{
		imports:                            importsHandler,
		instance:                           instance,
		proxyOnDownstreamData:              proxyOnDownstreamData,
		proxyOnUpstreamData:                proxyOnUpstreamData,
		proxyOnDownstreamConnectionClose:   proxyOnDownstreamConnectionClose,
		proxyOnUpstreamConnectionClose:     proxyOnUpstreamConnectionClose,
		proxyOnLog:                         proxyOnLog,
		proxyOnTick:                        proxyOnTick,
		proxyOnNewConnection:               proxyOnNewConnection,
		proxyOnContextCreate:               proxyOnContextCreate,
		proxyOnDone:                        proxyOnDone,
		proxyOnDelete:                      proxyOnDelete,
		proxyOnHttpCallResponse:            proxyOnHttpCallResponse,
		proxyOnQueueReady:                  proxyOnQueueReady,
		proxyOnGrpcClose:                   proxyOnGrpcClose,
		proxyOnGrpcReceiveInitialMetadata:  proxyOnGrpcReceiveInitialMetadata,
		proxyOnGrpcReceiveTrailingMetadata: proxyOnGrpcReceiveTrailingMetadata,
		proxyOnGrpcReceive:                 proxyOnGrpcReceive,
		proxyOnRequestBody:                 proxyOnRequestBody,
		proxyOnRequestHeaders:              proxyOnRequestHeaders,
		proxyOnRequestTrailers:             proxyOnRequestTrailers,
		proxyOnResponseBody:                proxyOnResponseBody,
		proxyOnResponseHeaders:             proxyOnResponseHeaders,
		proxyOnResponseTrailers:            proxyOnResponseTrailers,
	}, nil
}

type context struct {
	imports  api.ImportsHandler
	instance api.WasmInstance

	proxyOnDownstreamData              api.WasmFunction
	proxyOnUpstreamData                api.WasmFunction
	proxyOnDownstreamConnectionClose   api.WasmFunction
	proxyOnUpstreamConnectionClose     api.WasmFunction
	proxyOnLog                         api.WasmFunction
	proxyOnTick                        api.WasmFunction
	proxyOnNewConnection               api.WasmFunction
	proxyOnContextCreate               api.WasmFunction
	proxyOnDone                        api.WasmFunction
	proxyOnDelete                      api.WasmFunction
	proxyOnHttpCallResponse            api.WasmFunction
	proxyOnQueueReady                  api.WasmFunction
	proxyOnGrpcClose                   api.WasmFunction
	proxyOnGrpcReceiveInitialMetadata  api.WasmFunction
	proxyOnGrpcReceiveTrailingMetadata api.WasmFunction
	proxyOnGrpcReceive                 api.WasmFunction
	proxyOnRequestBody                 api.WasmFunction
	proxyOnRequestHeaders              api.WasmFunction
	proxyOnRequestTrailers             api.WasmFunction
	proxyOnResponseBody                api.WasmFunction
	proxyOnResponseHeaders             api.WasmFunction
	proxyOnResponseTrailers            api.WasmFunction
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

func (a *context) CallWasmFunction(funcName string, args ...interface{}) (interface{}, error) {
	ff, err := a.instance.GetExportsFunc(funcName)
	if err != nil {
		return api.ActionContinue, err
	}

	res, err := ff.Call(args...)
	if err != nil {
		a.instance.HandleError(err)

		return api.ActionContinue, err
	}

	return res, nil
}

// Configuration

func (a *context) ProxyOnVmStart(rootContextID int32, vmConfigurationSize int32) (bool, error) {
	res, err := a.CallWasmFunction("proxy_on_vm_start", rootContextID, vmConfigurationSize)
	if err != nil {
		return false, err
	}

	if v, ok := res.(int32); ok {
		return v == 1, nil
	}

	return false, ErrInvalidResult
}

func (a *context) ProxyOnConfigure(rootContextID int32, configurationSize int32) (bool, error) {
	res, err := a.CallWasmFunction("proxy_on_configure", rootContextID, configurationSize)
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
	_, err := a.proxyOnLog.Call(contextID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnTick(rootContextID int32) error {
	_, err := a.proxyOnTick.Call(rootContextID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnHttpCallResponse(contextID int32, tokenID int32, headerCount int32, bodySize int32, trailerCount int32) error {
	_, err := a.proxyOnHttpCallResponse.Call(contextID, tokenID, headerCount, bodySize, trailerCount)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnQueueReady(rootContextID int32, queueID int32) error {
	_, err := a.proxyOnQueueReady.Call(rootContextID, queueID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

// Context

func (a *context) ProxyOnContextCreate(contextID int32, rootContextID int32) error {
	_, err := a.proxyOnContextCreate.Call(contextID, rootContextID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnDone(contextID int32) (bool, error) {
	res, err := a.proxyOnDone.Call(contextID)
	if err != nil {
		a.instance.HandleError(err)
		return false, err
	}

	if v, ok := res.(int32); ok {
		return v == 1, nil
	}

	return false, ErrInvalidResult
}

func (a *context) ProxyOnDelete(contextID int32) error {
	_, err := a.proxyOnDelete.Call(contextID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

// L4

func (a *context) ProxyOnNewConnection(contextID int32) (api.Action, error) {
	action, err := a.proxyOnNewConnection.Call(contextID)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnDownstreamData(contextID int32, dataSize int32, endOfStream int32) (api.Action, error) {
	action, err := a.proxyOnDownstreamData.Call(contextID, dataSize, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnDownstreamConnectionClose(contextID int32, peerType int32) error {
	_, err := a.proxyOnDownstreamConnectionClose.Call(contextID, peerType)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnUpstreamData(contextID int32, dataSize int32, endOfStream int32) (api.Action, error) {
	action, err := a.proxyOnUpstreamData.Call(contextID, dataSize, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnUpstreamConnectionClose(contextID int32, peerType int32) error {
	_, err := a.proxyOnUpstreamConnectionClose.Call(contextID, peerType)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

// gRPC

func (a *context) ProxyOnGrpcClose(contextID int32, calloutID int32, statusCode int32) error {
	_, err := a.proxyOnGrpcClose.Call(contextID, calloutID, statusCode)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnGrpcReceiveInitialMetadata(contextID int32, tokenID int32, headerCount int32) error {
	_, err := a.proxyOnGrpcReceiveInitialMetadata.Call(contextID, tokenID, headerCount)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnGrpcReceiveTrailingMetadata(contextID int32, tokenID int32, trailerCount int32) error {
	_, err := a.proxyOnGrpcReceiveTrailingMetadata.Call(contextID, tokenID, trailerCount)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnGrpcReceive(contextID int32, tokenID int32, responseSize int32) error {
	_, err := a.proxyOnGrpcReceive.Call(contextID, tokenID, responseSize)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

// HTTP request

func (a *context) ProxyOnRequestBody(contextID int32, bodySize int32, endOfStream int32) (api.Action, error) {
	action, err := a.proxyOnRequestBody.Call(contextID, bodySize, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnRequestHeaders(contextID int32, headerCount int32, endOfStream int32) (api.Action, error) {
	action, err := a.proxyOnRequestHeaders.Call(contextID, headerCount, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnRequestTrailers(contextID int32, trailerCount int32) (api.Action, error) {
	action, err := a.proxyOnRequestTrailers.Call(contextID, trailerCount)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

// HTTP response

func (a *context) ProxyOnResponseBody(contextID int32, bodySize int32, endOfStream int32) (api.Action, error) {
	action, err := a.proxyOnResponseBody.Call(contextID, bodySize, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnResponseHeaders(contextID int32, headerCount int32, endOfStream int32) (api.Action, error) {
	action, err := a.proxyOnResponseHeaders.Call(contextID, headerCount, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnResponseTrailers(contextID int32, trailerCount int32) (api.Action, error) {
	action, err := a.proxyOnResponseTrailers.Call(contextID, trailerCount)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func unwrapAction(a interface{}) api.Action {
	if v, ok := a.(int32); ok {
		return api.Action(v)
	}

	return api.ActionContinue
}
