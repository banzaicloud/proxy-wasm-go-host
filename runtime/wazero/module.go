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
	"strings"

	"github.com/go-logr/logr"
	wazero "github.com/tetratelabs/wazero"
	klog "k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

type Module struct {
	vm       *VM
	module   wazero.CompiledModule
	rawBytes []byte
	ctx      context.Context
	logger   logr.Logger
}

func ModuleWithLogger(logger logr.Logger) ModuleOptions {
	return func(m *Module) {
		m.logger = logger
	}
}

type ModuleOptions func(module *Module)

func NewModule(ctx context.Context, vm *VM, module wazero.CompiledModule, wasmBytes []byte, options ...ModuleOptions) *Module {
	m := &Module{
		vm:       vm,
		module:   module,
		rawBytes: wasmBytes,
		ctx:      ctx,
	}

	for _, option := range options {
		option(m)
	}

	if vm.logger == (logr.Logger{}) {
		vm.logger = klog.Background()
	}

	return m
}

func (m *Module) Init() {}

func (m *Module) NewInstance() (api.WasmInstance, error) {
	return NewInstance(m.ctx, m.vm, m, InstanceWithLogger(m.logger)), nil
}

func (m *Module) GetABINameList() []string {
	abiNameList := make([]string, 0)

	exportList := m.module.ExportedFunctions()

	for export := range exportList {
		if strings.HasPrefix(export, "proxy_abi") {
			abiNameList = append(abiNameList, export)
		}
	}

	return abiNameList
}
