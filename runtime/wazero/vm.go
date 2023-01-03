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

package wazero

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	wazero "github.com/tetratelabs/wazero"
	klog "k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

type VM struct {
	runtime wazero.Runtime
	ctx     context.Context
	logger  logr.Logger
}

func VMWithLogger(logger logr.Logger) VMOptions {
	return func(vm *VM) {
		vm.logger = logger
	}
}

type VMOptions func(vm *VM)

func NewVM(ctx context.Context, options ...VMOptions) api.WasmVM {
	vm := &VM{
		runtime: wazero.NewRuntime(ctx),
		ctx:     ctx,
	}

	for _, option := range options {
		option(vm)
	}

	if vm.logger == (logr.Logger{}) {
		vm.logger = klog.Background()
	}

	return vm
}

func (w *VM) Name() string {
	return "wazero"
}

func (w *VM) Init() {}

func (w *VM) NewModule(wasmBytes []byte) (api.WasmModule, error) {
	if len(wasmBytes) == 0 {
		return nil, errors.New("wasm was empty")
	}

	m, err := w.runtime.CompileModule(w.ctx, wasmBytes)
	if err != nil {
		return nil, errors.WrapIf(err, "could not compile module")
	}

	return NewModule(w.ctx, w, m, wasmBytes, ModuleWithLogger(w.logger)), nil
}

// Close implements io.Closer
func (w *VM) Close() (err error) {
	if r := w.runtime; r != nil {
		err = r.Close(w.ctx)
	}
	return
}
