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

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/abi"
	"github.com/banzaicloud/proxy-wasm-go-host/api"
	"github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmer"
)

var (
	contextIDGenerator int32
	rootContextID      int32
)

var (
	lock    sync.Mutex
	once    sync.Once
	wasmCtx api.ABIContext
)

var _ api.ImportsHandler = &importHandler{}

// implement v1.ImportsHandler.
type importHandler struct {
	reqHeader api.HeaderMap
	abi.DefaultImportsHandler
}

// override.
func (im *importHandler) GetHttpRequestHeader() api.HeaderMap {
	return im.reqHeader
}

// override.
func (im *importHandler) Log(level api.LogLevel, msg string) api.WasmResult {
	fmt.Println(msg)
	return api.WasmResultOk
}

// serve HTTP req
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("receive request %s\n", r.URL)
	for k, v := range r.Header {
		fmt.Printf("print header from server host, %v -> %v\n", k, v)
	}

	// get wasm vm instance
	ctx := getWasmContext()

	// create context id for the http req
	contextID := atomic.AddInt32(&contextIDGenerator, 1)

	// do wasm

	// according to ABI, we should create a root context id before any operations
	once.Do(func() {
		if err := ctx.GetExports().ProxyOnContextCreate(rootContextID, 0); err != nil {
			log.Panicln(err)
		}
	})

	// lock wasm vm instance for exclusive ownership
	ctx.GetInstance().Lock(ctx)
	defer ctx.GetInstance().Unlock()

	// Set the import handler to the current request.
	ctx.SetImports(&importHandler{reqHeader: &myHeaderMap{r.Header}})

	// create wasm-side context id for current http req
	if err := ctx.GetExports().ProxyOnContextCreate(contextID, rootContextID); err != nil {
		log.Panicln(err)
	}

	// call wasm-side on_request_header
	if _, err := ctx.GetExports().ProxyOnRequestHeaders(contextID, int32(len(r.Header)), 1); err != nil {
		log.Panicln(err)
	}

	// delete wasm-side context id to prevent memory leak
	if err := ctx.GetExports().ProxyOnDelete(contextID); err != nil {
		log.Panicln(err)
	}

	// reply with ok
	w.WriteHeader(http.StatusOK)
}

var vm api.WasmVM

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	defer func() {
		if vm != nil {
			vm.Close()
		}
	}()

	// create root context id
	rootContextID = atomic.AddInt32(&contextIDGenerator, 1)

	// serve http
	http.HandleFunc("/", ServeHTTP)

	server := &http.Server{
		Addr:              "127.0.0.1:2045",
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panicln(err)
	}
}

func getWasmContext() api.ABIContext {
	lock.Lock()
	defer lock.Unlock()

	if wasmCtx == nil {
		guest, err := os.ReadFile("data/http.wasm")
		if err != nil {
			log.Panicln(err)
		}

		vm = wasmer.NewWasmerVM(wasmer.VMWithLogger(klog.Background()))

		module, err := vm.NewModule(guest)
		if err != nil {
			log.Panicln(err)
		}

		instance, err := module.NewInstance()
		if err != nil {
			log.Panicln(err)
		}

		// create ABI context
		wasmCtx = abi.NewContext(&abi.DefaultImportsHandler{}, instance)

		// start the wasm vm instance
		if err = instance.Start(); err != nil {
			log.Panicln(err)
		}
	}

	return wasmCtx
}

// wrapper for http.Header, convert Header to api.HeaderMap.
type myHeaderMap struct {
	realMap http.Header
}

func (m *myHeaderMap) Get(key string) (string, bool) {
	return m.realMap.Get(key), true
}

func (m *myHeaderMap) Set(key, value string) { panic("implemented") }

func (m *myHeaderMap) Add(key, value string) { panic("implemented") }

func (m *myHeaderMap) Del(key string) { panic("implemented") }

func (m *myHeaderMap) Range(f func(key string, value string) bool) {
	for k := range m.realMap {
		v := m.realMap.Get(k)
		f(k, v)
	}
}

func (m *myHeaderMap) Clone() api.HeaderMap { panic("implemented") }

func (m *myHeaderMap) ByteSize() uint64 { panic("implemented") }
