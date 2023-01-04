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
	"strings"

	"github.com/go-logr/logr"
	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
	"k8s.io/klog/v2"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

type Module struct {
	vm          *VM
	module      *wasmerGo.Module
	abiNameList []string
	wasiVersion wasmerGo.WasiVersion
	debug       *dwarfInfo
	rawBytes    []byte
	logger      logr.Logger
}

func ModuleWithLogger(logger logr.Logger) ModuleOptions {
	return func(m *Module) {
		m.logger = logger
	}
}

type ModuleOptions func(module *Module)

func NewWasmerModule(vm *VM, module *wasmerGo.Module, wasmBytes []byte, options ...ModuleOptions) *Module {
	m := &Module{
		vm:       vm,
		module:   module,
		rawBytes: wasmBytes,
	}

	for _, option := range options {
		option(m)
	}

	if vm.logger == (logr.Logger{}) {
		vm.logger = klog.Background()
	}

	m.Init()

	return m
}

func (m *Module) Init() {
	m.wasiVersion = wasmerGo.GetWasiVersion(m.module)

	m.abiNameList = m.GetABINameList()

	// parse dwarf info from wasm data bytes
	if debug := parseDwarf(m.rawBytes); debug != nil {
		m.debug = debug
	}

	// release raw bytes, the parsing of dwarf info is the only place that uses module raw bytes
	m.rawBytes = nil
}

func (m *Module) NewInstance() (api.WasmInstance, error) {
	if m.debug != nil {
		return NewWasmerInstance(m.vm, m, InstanceWithDebug(m.debug), InstanceWithLogger(m.logger))
	}

	return NewWasmerInstance(m.vm, m)
}

func (m *Module) GetABINameList() []string {
	abiNameList := make([]string, 0)

	exportList := m.module.Exports()

	for _, export := range exportList {
		if export.Type().Kind() == wasmerGo.FUNCTION {
			if strings.HasPrefix(export.Name(), "proxy_abi") {
				abiNameList = append(abiNameList, export.Name())
			}
		}
	}

	return abiNameList
}
