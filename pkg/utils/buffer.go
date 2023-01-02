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

// CommonBuffer is a simple implementation of IoBuffer.
type CommonBuffer struct {
	buf []byte
}

func (c *CommonBuffer) Len() int {
	return len(c.buf)
}

func (c *CommonBuffer) Bytes() []byte {
	return c.buf
}

func (c *CommonBuffer) Write(p []byte) (int, error) {
	c.buf = append(c.buf, p...)

	return len(p), nil
}

func (c *CommonBuffer) Drain(offset int) {
	if offset > len(c.buf) {
		return
	}
	c.buf = c.buf[offset:]
}

func NewIoBufferBytes(data []byte) api.IoBuffer {
	return &CommonBuffer{buf: data}
}
