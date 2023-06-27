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
	"emperror.dev/errors"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

const ProxyWasmABI_0_1_0 string = "proxy_abi_version_0_1_0"
const ProxyWasmABI_0_2_0 string = "proxy_abi_version_0_2_0"
const ProxyWasmABI_0_2_1 string = "proxy_abi_version_0_2_1"

var ErrFuncNotExported = errors.New("function is not exported")

func NewContext(importsHandler api.ImportsHandler, instance api.WasmInstance) api.ABIContext {
	getWasmFunction := func(name string) (f api.WasmFunction) {
		f, _ = instance.GetExportsFunc(name)

		return
	}

	return &context{
		imports:                            importsHandler,
		instance:                           instance,
		proxyOnVmStart:                     getWasmFunction("proxy_on_vm_start"),
		proxyOnConfigure:                   getWasmFunction("proxy_on_configure"),
		proxyOnDownstreamData:              getWasmFunction("proxy_on_downstream_data"),
		proxyOnUpstreamData:                getWasmFunction("proxy_on_upstream_data"),
		proxyOnDownstreamConnectionClose:   getWasmFunction("proxy_on_upstream_connection_close"),
		proxyOnUpstreamConnectionClose:     getWasmFunction("proxy_on_downstream_connection_close"),
		proxyOnLog:                         getWasmFunction("proxy_on_log"),
		proxyOnTick:                        getWasmFunction("proxy_on_tick"),
		proxyOnNewConnection:               getWasmFunction("proxy_on_new_connection"),
		proxyOnContextCreate:               getWasmFunction("proxy_on_context_create"),
		proxyOnDone:                        getWasmFunction("proxy_on_done"),
		proxyOnDelete:                      getWasmFunction("proxy_on_delete"),
		proxyOnHttpCallResponse:            getWasmFunction("proxy_on_http_call_response"),
		proxyOnQueueReady:                  getWasmFunction("proxy_on_queue_ready"),
		proxyOnGrpcClose:                   getWasmFunction("proxy_on_grpc_close"),
		proxyOnGrpcReceiveInitialMetadata:  getWasmFunction("proxy_on_grpc_receive_initial_metadata"),
		proxyOnGrpcReceiveTrailingMetadata: getWasmFunction("proxy_on_grpc_receive_trailing_metadata"),
		proxyOnGrpcReceive:                 getWasmFunction("proxy_on_grpc_receive"),
		proxyOnRequestBody:                 getWasmFunction("proxy_on_request_body"),
		proxyOnRequestHeaders:              getWasmFunction("proxy_on_request_headers"),
		proxyOnRequestTrailers:             getWasmFunction("proxy_on_request_trailers"),
		proxyOnResponseBody:                getWasmFunction("proxy_on_response_body"),
		proxyOnResponseHeaders:             getWasmFunction("proxy_on_response_headers"),
		proxyOnResponseTrailers:            getWasmFunction("proxy_on_response_trailers"),
	}
}

type context struct {
	imports  api.ImportsHandler
	instance api.WasmInstance

	proxyOnVmStart                     api.WasmFunction
	proxyOnConfigure                   api.WasmFunction
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
	if a.proxyOnVmStart == nil {
		return false, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_vm_start")
	}

	res, err := a.proxyOnVmStart.Call(rootContextID, vmConfigurationSize)
	if err != nil {
		return false, err
	}

	if v, ok := res.(int32); ok {
		return v == 1, nil
	}

	return false, ErrInvalidResult
}

func (a *context) ProxyOnConfigure(rootContextID int32, configurationSize int32) (bool, error) {
	if a.proxyOnConfigure == nil {
		return false, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_configure")
	}

	res, err := a.proxyOnConfigure.Call(rootContextID, configurationSize)
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
	if a.proxyOnLog == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_log")
	}

	_, err := a.proxyOnLog.Call(contextID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnTick(rootContextID int32) error {
	if a.proxyOnTick == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_tick")
	}

	_, err := a.proxyOnTick.Call(rootContextID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnHttpCallResponse(contextID int32, tokenID int32, headerCount int32, bodySize int32, trailerCount int32) error {
	if a.proxyOnHttpCallResponse == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_http_call_response")
	}

	_, err := a.proxyOnHttpCallResponse.Call(contextID, tokenID, headerCount, bodySize, trailerCount)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnQueueReady(rootContextID int32, queueID int32) error {
	if a.proxyOnQueueReady == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_queue_ready")
	}

	_, err := a.proxyOnQueueReady.Call(rootContextID, queueID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

// Context

func (a *context) ProxyOnContextCreate(contextID int32, rootContextID int32) error {
	if a.proxyOnContextCreate == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_context_create")
	}

	_, err := a.proxyOnContextCreate.Call(contextID, rootContextID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnDone(contextID int32) (bool, error) {
	if a.proxyOnDone == nil {
		return false, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_done")
	}

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
	if a.proxyOnDelete == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_delete")
	}

	_, err := a.proxyOnDelete.Call(contextID)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

// L4

func (a *context) ProxyOnNewConnection(contextID int32) (api.Action, error) {
	if a.proxyOnNewConnection == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_new_connection")
	}

	action, err := a.proxyOnNewConnection.Call(contextID)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnDownstreamData(contextID int32, dataSize int32, endOfStream int32) (api.Action, error) {
	if a.proxyOnDownstreamData == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_downstream_data")
	}

	action, err := a.proxyOnDownstreamData.Call(contextID, dataSize, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnDownstreamConnectionClose(contextID int32, peerType int32) error {
	if a.proxyOnDownstreamConnectionClose == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_downstream_connection_close")
	}

	_, err := a.proxyOnDownstreamConnectionClose.Call(contextID, peerType)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnUpstreamData(contextID int32, dataSize int32, endOfStream int32) (api.Action, error) {
	if a.proxyOnUpstreamData == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_upstream_data")
	}

	action, err := a.proxyOnUpstreamData.Call(contextID, dataSize, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnUpstreamConnectionClose(contextID int32, peerType int32) error {
	if a.proxyOnUpstreamConnectionClose == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_upstream_connection_close")
	}

	_, err := a.proxyOnUpstreamConnectionClose.Call(contextID, peerType)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

// gRPC

func (a *context) ProxyOnGrpcClose(contextID int32, calloutID int32, statusCode int32) error {
	if a.proxyOnGrpcClose == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_grpc_close")
	}

	_, err := a.proxyOnGrpcClose.Call(contextID, calloutID, statusCode)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnGrpcReceiveInitialMetadata(contextID int32, tokenID int32, headerCount int32) error {
	if a.proxyOnGrpcReceiveInitialMetadata == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_grpc_receive_initial_metadata")
	}

	_, err := a.proxyOnGrpcReceiveInitialMetadata.Call(contextID, tokenID, headerCount)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnGrpcReceiveTrailingMetadata(contextID int32, tokenID int32, trailerCount int32) error {
	if a.proxyOnGrpcReceiveTrailingMetadata == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_grpc_receive_trailing_metadata")
	}

	_, err := a.proxyOnGrpcReceiveTrailingMetadata.Call(contextID, tokenID, trailerCount)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

func (a *context) ProxyOnGrpcReceive(contextID int32, tokenID int32, responseSize int32) error {
	if a.proxyOnGrpcReceive == nil {
		return errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_grpc_receive")
	}

	_, err := a.proxyOnGrpcReceive.Call(contextID, tokenID, responseSize)
	if err != nil {
		a.instance.HandleError(err)
	}

	return err
}

// HTTP request

func (a *context) ProxyOnRequestBody(contextID int32, bodySize int32, endOfStream int32) (api.Action, error) {
	if a.proxyOnRequestBody == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_request_body")
	}

	action, err := a.proxyOnRequestBody.Call(contextID, bodySize, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnRequestHeaders(contextID int32, headerCount int32, endOfStream int32) (api.Action, error) {
	if a.proxyOnRequestHeaders == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_request_headers")
	}

	action, err := a.proxyOnRequestHeaders.Call(contextID, headerCount, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnRequestTrailers(contextID int32, trailerCount int32) (api.Action, error) {
	if a.proxyOnRequestTrailers == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_request_trailers")
	}

	action, err := a.proxyOnRequestTrailers.Call(contextID, trailerCount)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

// HTTP response

func (a *context) ProxyOnResponseBody(contextID int32, bodySize int32, endOfStream int32) (api.Action, error) {
	if a.proxyOnResponseBody == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_response_body")
	}

	action, err := a.proxyOnResponseBody.Call(contextID, bodySize, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnResponseHeaders(contextID int32, headerCount int32, endOfStream int32) (api.Action, error) {
	if a.proxyOnResponseHeaders == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_response_headers")
	}

	action, err := a.proxyOnResponseHeaders.Call(contextID, headerCount, endOfStream)
	if err != nil {
		a.instance.HandleError(err)
		return api.ActionPause, err
	}

	return unwrapAction(action), nil
}

func (a *context) ProxyOnResponseTrailers(contextID int32, trailerCount int32) (api.Action, error) {
	if a.proxyOnResponseTrailers == nil {
		return api.ActionContinue, errors.WithDetails(ErrFuncNotExported, "func", "proxy_on_response_trailers")
	}

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
