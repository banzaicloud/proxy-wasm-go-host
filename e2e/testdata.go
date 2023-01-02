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
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wabin/binary"
	"github.com/tetratelabs/wabin/wasm"
)

var binAddRequestHeader = func() []byte {
	return loadWasmWithABI("testdata/add-req-header/main.wasm", "proxy_abi_version_0_2_0")
}()

func loadWasmWithABI(wasmPath, abiName string) []byte {
	bin, err := os.ReadFile(wasmPath)
	if err != nil {
		log.Panicln(err)
	}
	mod, err := binary.DecodeModule(bin, wasm.CoreFeaturesV2)
	if err != nil {
		log.Panicln(err)
	}
	exports := []string{}
	for _, e := range mod.ExportSection {
		if e.Name == abiName {
			return bin
		}
		exports = append(exports, e.Name)
	}
	log.Panicln(fmt.Errorf("export not found in %s, want: %v, have: %v", wasmPath, abiName, exports))
	return nil
}
