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

package wasmer

import (
	"emperror.dev/errors"
	"github.com/go-logr/logr"
	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
	"k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

type VM struct {
	engine *wasmerGo.Engine
	store  *wasmerGo.Store
	logger logr.Logger
}

func VMWithLogger(logger logr.Logger) VMOptions {
	return func(vm *VM) {
		vm.logger = logger
	}
}

type VMOptions func(vm *VM)

func NewWasmerVM(options ...VMOptions) api.WasmVM {
	vm := &VM{}
	vm.Init()

	for _, option := range options {
		option(vm)
	}

	if vm.logger == (logr.Logger{}) {
		vm.logger = klog.Background()
	}

	return vm
}

func (w *VM) Name() string {
	return "wasmer"
}

func (w *VM) Init() {
	w.engine = wasmerGo.NewEngine()
	w.store = wasmerGo.NewStore(w.engine)
}

func (w *VM) NewModule(wasmBytes []byte) (api.WasmModule, error) {
	if len(wasmBytes) == 0 {
		return nil, errors.New("wasm was empty")
	}

	m, err := wasmerGo.NewModule(w.store, wasmBytes)
	if err != nil {
		return nil, errors.WrapIf(err, "could not instantiate module")
	}

	return NewWasmerModule(w, m, wasmBytes, ModuleWithLogger(w.logger)), nil
}

// Close implements io.Closer
func (w *VM) Close() (err error) {
	if s := w.store; s != nil {
		s.Close()
	}
	return
}
