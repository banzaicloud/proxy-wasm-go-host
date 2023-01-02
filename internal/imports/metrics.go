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

// Metrics

func (h *host) ProxyDefineMetric(ctx context.Context, metricType int32, namePtr int32, nameSize int32, returnMetricId int32) int32 {
	instance := h.Instance
	ih := getImportHandler(instance)

	if api.MetricType(metricType) > api.MetricTypeMax {
		return api.WasmResultBadArgument.Int32()
	}

	name, err := instance.GetMemory(uint64(namePtr), uint64(nameSize))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}
	if len(name) == 0 {
		return api.WasmResultBadArgument.Int32()
	}

	mid, res := ih.DefineMetric(api.MetricType(metricType), string(name))
	if res != api.WasmResultOk {
		return res.Int32()
	}

	err = instance.PutUint32(uint64(returnMetricId), uint32(mid))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}

func (h *host) ProxyIncrementMetric(ctx context.Context, metricId int32, offset int64) int32 {
	ih := getImportHandler(h.Instance)

	res := ih.IncrementMetric(metricId, offset)

	return res.Int32()
}

func (h *host) ProxyRecordMetric(ctx context.Context, metricId int32, value int64) int32 {
	ih := getImportHandler(h.Instance)

	res := ih.RecordMetric(metricId, value)

	return res.Int32()
}

func (h *host) ProxyGetMetric(ctx context.Context, metricId int32, resultUint64Ptr int32) int32 {
	instance := h.Instance
	ih := getImportHandler(instance)

	value, res := ih.GetMetric(metricId)
	if res != api.WasmResultOk {
		return res.Int32()
	}

	err := instance.PutUint32(uint64(resultUint64Ptr), uint32(value))
	if err != nil {
		return api.WasmResultInvalidMemoryAccess.Int32()
	}

	return api.WasmResultOk.Int32()
}