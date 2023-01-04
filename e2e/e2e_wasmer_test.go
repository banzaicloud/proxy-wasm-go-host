//go:build wasmer
// +build wasmer

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

package e2e

import (
	_ "embed"
	"testing"

	"github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmer"
)

func TestStartABIContextV1_wasmer(t *testing.T) {
	vm := wasmer.NewWasmerVM()
	defer vm.Close()

	testStartABIContext(t, vm)
}

func TestAddRequestHeaderV1_wasmer(t *testing.T) {
	vm := wasmer.NewWasmerVM()
	defer vm.Close()

	testV1(t, vm, testAddRequestHeader)
}
