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

package utils

import "github.com/banzaicloud/proxy-wasm-go-host/api"

// CommonHeader is a simple implementation of HeaderMap.
type CommonHeader map[string]string

func (h CommonHeader) Get(key string) (value string, ok bool) {
	value, ok = h[key]

	return
}

func (h CommonHeader) Set(key string, value string) {
	h[key] = value
}

func (h CommonHeader) Add(key string, value string) {
	panic("not supported")
}

func (h CommonHeader) Del(key string) {
	delete(h, key)
}

func (h CommonHeader) Range(f func(key, value string) bool) {
	for k, v := range h {
		// stop if f return false
		if !f(k, v) {
			break
		}
	}
}

func (h CommonHeader) Clone() api.HeaderMap {
	c := make(map[string]string)

	for k, v := range h {
		c[k] = v
	}

	return CommonHeader(c)
}

func (h CommonHeader) ByteSize() uint64 {
	var size uint64

	for k, v := range h {
		size += uint64(len(k) + len(v))
	}

	return size
}