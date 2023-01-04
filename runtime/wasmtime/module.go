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

package wasmtime

import (
	"context"
	"strings"

	"github.com/bytecodealliance/wasmtime-go/v3"
	"github.com/go-logr/logr"
	klog "k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

type Module struct {
	vm       *VM
	module   *wasmtime.Module
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

func NewModule(ctx context.Context, vm *VM, module *wasmtime.Module, wasmBytes []byte, options ...ModuleOptions) *Module {
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

	exportList := m.module.Exports()

	for _, export := range exportList {
		if strings.HasPrefix(export.Name(), "proxy_abi") {
			abiNameList = append(abiNameList, export.Name())
		}
	}

	return abiNameList
}
