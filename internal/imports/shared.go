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

package imports

import (
	"context"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

// SharedData

func (h *host) ProxyGetSharedData(ctx context.Context, keyPtr int32, keySize int32, returnValuePtr int32, returnValueSize int32, returnCasPtr int32) int32 {
	instance := h.Instance
	key, err := instance.GetMemory(uint64(keyPtr), uint64(keySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	ih := getImportHandler(instance)

	v, cas, res := ih.GetSharedData(string(key))
	if res != api.WasmResultOk {
		return res.Int32()
	}

	res = copyIntoInstance(instance, v, returnValuePtr, returnValueSize)
	if res != api.WasmResultOk {
		return res.Int32()
	}

	err = instance.PutUint32(uint64(returnCasPtr), cas)
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

func (h *host) ProxySetSharedData(ctx context.Context, keyPtr int32, keySize int32, valuePtr int32, valueSize int32, cas int32) int32 {
	instance := h.Instance
	key, err := instance.GetMemory(uint64(keyPtr), uint64(keySize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(key) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	value, err := instance.GetMemory(uint64(valuePtr), uint64(valueSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	ih := getImportHandler(instance)

	return ih.SetSharedData(string(key), string(value), uint32(cas)).Int32()
}

// SharedQueue

func (h *host) ProxyRegisterSharedQueue(ctx context.Context, queueNamePtr int32, queueNameSize int32, tokenIDPtr int32) int32 {
	instance := h.Instance
	queueName, err := instance.GetMemory(uint64(queueNamePtr), uint64(queueNameSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(queueName) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	ih := getImportHandler(instance)

	queueID, res := ih.RegisterSharedQueue(string(queueName))
	if res != api.WasmResultOk {
		return res.Int32()
	}

	err = instance.PutUint32(uint64(tokenIDPtr), queueID)
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

// TODO(@wayne): vm context
func (h *host) ProxyResolveSharedQueue(ctx context.Context, vmIDPtr int32, vmIDSize int32, queueNamePtr int32, queueNameSize int32, tokenIDPtr int32) int32 {
	instance := h.Instance
	queueName, err := instance.GetMemory(uint64(queueNamePtr), uint64(queueNameSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(queueName) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	ih := getImportHandler(instance)

	queueID, res := ih.ResolveSharedQueue(string(queueName))
	if res != api.WasmResultOk {
		return res.Int32()
	}

	err = instance.PutUint32(uint64(tokenIDPtr), queueID)
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

func (h *host) ProxyEnqueueSharedQueue(ctx context.Context, tokenID int32, dataPtr int32, dataSize int32) int32 {
	instance := h.Instance
	value, err := instance.GetMemory(uint64(dataPtr), uint64(dataSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	ih := getImportHandler(instance)

	return ih.EnqueueSharedQueue(uint32(tokenID), string(value)).Int32()
}

func (h *host) ProxyDequeueSharedQueue(ctx context.Context, tokenID int32, returnValuePtr int32, returnValueSize int32) int32 {
	instance := h.Instance
	ih := getImportHandler(instance)

	value, res := ih.DequeueSharedQueue(uint32(tokenID))
	if res != api.WasmResultOk {
		return res.Int32()
	}

	return copyIntoInstance(instance, value, returnValuePtr, returnValueSize).Int32()
}
