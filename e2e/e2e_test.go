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

//nolint:goerr113
package e2e

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"testing"

	"mosn.io/pkg/log"

	"github.com/banzaicloud/proxy-wasm-go-host/abi"
	"github.com/banzaicloud/proxy-wasm-go-host/api"
	"github.com/banzaicloud/proxy-wasm-go-host/pkg/utils"
	"github.com/banzaicloud/proxy-wasm-go-host/runtime/wazero"
)

func init() {
	log.DefaultLogger.SetLogLevel(log.ERROR)
}

func TestStartABIContext_wazero(t *testing.T) {
	t.Parallel()

	vm := wazero.NewVM(context.Background())
	defer vm.Close()

	testStartABIContext(t, vm)
}

func testStartABIContext(t *testing.T, vm api.WasmVM) {
	t.Helper()

	module, err := vm.NewModule(binAddRequestHeader)
	if err != nil {
		t.Fatal(err)
	}

	instance, err := module.NewInstance()
	if err != nil {
		t.Fatal(err)
	}

	defer instance.Stop()

	if _, err := startABIContext(instance); err != nil {
		t.Fatal(err)
	}
}

func startABIContext(instance api.WasmInstance) (wasmCtx api.ABIContext, err error) {
	// create ABI context
	wasmCtx = abi.NewContext(&abi.DefaultImportsHandler{}, instance)

	// start the wasm vm instance
	err = instance.Start()
	return
}

func TestAddRequestHeader_wazero(t *testing.T) {
	t.Parallel()

	vm := wazero.NewVM(context.Background())
	defer vm.Close()

	testV1(t, vm, testAddRequestHeader)
}

func testAddRequestHeader(wasmCtx api.ABIContext, contextID int32) error {
	handler := &headersHandler{reqHeader: &utils.CommonHeader{}}
	wasmCtx.SetImports(handler)

	if action, err := wasmCtx.GetExports().ProxyOnRequestHeaders(contextID, 0, 1); err != nil {
		return err
	} else if want, have := api.ActionContinue, action; want != have {
		return fmt.Errorf("unexpected action, want: %v, have: %v", want, have)
	}

	expectedHeader := "Wasm-Context"
	want := strconv.Itoa(int(contextID))
	if have, _ := handler.GetHttpRequestHeader().Get(expectedHeader); want != have {
		return fmt.Errorf("unexpected %s, want: %v, have: %v", expectedHeader, want, have)
	}
	return nil
}

func testV1(t *testing.T, vm api.WasmVM, test func(wasmCtx api.ABIContext, contextID int32) error) {
	t.Helper()

	module, err := vm.NewModule(binAddRequestHeader)
	if err != nil {
		t.Fatal(err)
	}

	instance, err := module.NewInstance()
	if err != nil {
		t.Fatal(err)
	}

	defer instance.Stop()

	wasmCtx, err := startABIContext(instance)
	if err != nil {
		t.Fatal(err)
	}
	defer wasmCtx.GetInstance().Stop()

	exports := wasmCtx.GetExports()

	// make the root context
	rootContextID := int32(1)
	if err := exports.ProxyOnContextCreate(rootContextID, int32(0)); err != nil {
		t.Fatal(err)
	}

	// lock wasm vm instance for exclusive ownership
	wasmCtx.GetInstance().Lock(wasmCtx)
	defer wasmCtx.GetInstance().Unlock()

	contextID := int32(2)
	if err = exports.ProxyOnContextCreate(contextID, rootContextID); err != nil {
		t.Fatal(err)
	}

	if err = test(wasmCtx, contextID); err != nil {
		t.Fatal(err)
	}

	if _, err = exports.ProxyOnDone(contextID); err != nil {
		t.Fatal(err)
	}

	if err = exports.ProxyOnDelete(contextID); err != nil {
		t.Fatal(err)
	}
}

var _ api.ImportsHandler = &headersHandler{}

// headersHandler implements api.ImportsHandler.
type headersHandler struct {
	reqHeader api.HeaderMap
	abi.DefaultImportsHandler
}

// override.
func (im *headersHandler) GetHttpRequestHeader() api.HeaderMap {
	return im.reqHeader
}
