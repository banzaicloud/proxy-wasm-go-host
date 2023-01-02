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

package abi

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/banzaicloud/proxy-wasm-go-host/api"
)

type DefaultImportsHandler struct{}

// for golang host environment, no-op
func (d *DefaultImportsHandler) Wait() api.Action { return api.ActionContinue }

// utils
func (d *DefaultImportsHandler) GetRootContextID() int32 { return 0 }

func (d *DefaultImportsHandler) GetVmConfig() api.IoBuffer { return nil }

func (d *DefaultImportsHandler) GetPluginConfig() api.IoBuffer { return nil }

func (d *DefaultImportsHandler) Log(level api.LogLevel, msg string) api.WasmResult {
	return api.WasmResultOk
}

func (d *DefaultImportsHandler) GetLogLevel() api.LogLevel {
	return api.LogLevelInfo
}

func (d *DefaultImportsHandler) GetStatus() (int32, string, api.WasmResult) {
	return http.StatusOK, http.StatusText(http.StatusOK), api.WasmResultOk
}

func (d *DefaultImportsHandler) SetEffectiveContextID(contextID int32) api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) SetTickPeriodMilliseconds(tickPeriodMilliseconds int32) api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) GetCurrentTimeNanoseconds() (int32, api.WasmResult) {
	nano := time.Now().Nanosecond()

	return int32(nano), api.WasmResultOk
}

func (d *DefaultImportsHandler) Done() api.WasmResult { return api.WasmResultUnimplemented }

// l4

func (d *DefaultImportsHandler) GetDownStreamData() api.IoBuffer { return nil }

func (d *DefaultImportsHandler) GetUpstreamData() api.IoBuffer { return nil }

func (d *DefaultImportsHandler) ResumeDownstream() api.WasmResult { return api.WasmResultUnimplemented }

func (d *DefaultImportsHandler) ResumeUpstream() api.WasmResult { return api.WasmResultUnimplemented }

func (d *DefaultImportsHandler) CloseDownstream() api.WasmResult { return api.WasmResultUnimplemented }

func (d *DefaultImportsHandler) CloseUpstream() api.WasmResult { return api.WasmResultUnimplemented }

// http

func (d *DefaultImportsHandler) GetHttpRequestHeader() api.HeaderMap { return nil }

func (d *DefaultImportsHandler) GetHttpRequestBody() api.IoBuffer { return nil }

func (d *DefaultImportsHandler) GetHttpRequestTrailer() api.HeaderMap { return nil }

func (d *DefaultImportsHandler) GetHttpResponseHeader() api.HeaderMap { return nil }

func (d *DefaultImportsHandler) GetHttpResponseBody() api.IoBuffer { return nil }

func (d *DefaultImportsHandler) GetHttpResponseTrailer() api.HeaderMap { return nil }

func (d *DefaultImportsHandler) GetHttpCallResponseHeaders() api.HeaderMap { return nil }

func (d *DefaultImportsHandler) GetHttpCallResponseBody() api.IoBuffer { return nil }

func (d *DefaultImportsHandler) GetHttpCallResponseTrailer() api.HeaderMap { return nil }

func (d *DefaultImportsHandler) HttpCall(url string, headers api.HeaderMap, body api.IoBuffer, trailer api.HeaderMap, timeoutMilliseconds int32) (int32, api.WasmResult) {
	return 0, api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) ResumeHttpRequest() api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) ResumeHttpResponse() api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) CloseHttpRequest() api.WasmResult { return api.WasmResultUnimplemented }

func (d *DefaultImportsHandler) CloseHttpResponse() api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) SendHttpResp(respCode int32, respCodeDetail api.IoBuffer, respBody api.IoBuffer, additionalHeaderMap api.HeaderMap, grpcCode int32) api.WasmResult {
	return api.WasmResultUnimplemented
}

// grpc

func (d *DefaultImportsHandler) OpenGrpcStream(grpcService string, serviceName string, method string) (int32, api.WasmResult) {
	return 0, api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) SendGrpcCallMsg(token int32, data api.IoBuffer, endOfStream int32) api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) CancelGrpcCall(token int32) api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) CloseGrpcCall(token int32) api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) GrpcCall(grpcService string, serviceName string, method string, data api.IoBuffer, timeoutMilliseconds int32) (int32, api.WasmResult) {
	return 0, api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) GetGrpcReceiveInitialMetaData() api.HeaderMap { return nil }

func (d *DefaultImportsHandler) GetGrpcReceiveBuffer() api.IoBuffer { return nil }

func (d *DefaultImportsHandler) GetGrpcReceiveTrailerMetaData() api.HeaderMap { return nil }

// foreign

func (d *DefaultImportsHandler) CallForeignFunction(funcName string, param []byte) ([]byte, api.WasmResult) {
	return nil, api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) GetFuncCallData() api.IoBuffer { return nil }

// property

func (d *DefaultImportsHandler) GetProperty(key string) (string, api.WasmResult) {
	return "", api.WasmResultOk
}

func (d *DefaultImportsHandler) SetProperty(key string, value string) api.WasmResult {
	return api.WasmResultUnimplemented
}

// metric

func (d *DefaultImportsHandler) DefineMetric(metricType api.MetricType, name string) (int32, api.WasmResult) {
	return 0, api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) IncrementMetric(metricID int32, offset int64) api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) RecordMetric(metricID int32, value int64) api.WasmResult {
	return api.WasmResultUnimplemented
}

func (d *DefaultImportsHandler) GetMetric(metricID int32) (int64, api.WasmResult) {
	return 0, api.WasmResultUnimplemented
}

// shared data

type sharedDataItem struct {
	data string
	cas  uint32
}

type sharedData struct {
	lock sync.RWMutex
	m    map[string]*sharedDataItem
	cas  uint32
}

func (s *sharedData) get(key string) (string, uint32, api.WasmResult) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if v, ok := s.m[key]; ok {
		return v.data, v.cas, api.WasmResultOk
	}

	return "", 0, api.WasmResultNotFound
}

func (s *sharedData) set(key string, value string, cas uint32) api.WasmResult {
	if key == "" {
		return api.WasmResultBadArgument
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.m[key]; ok {
		if v.cas != cas {
			return api.WasmResultCasMismatch
		}

		v.data = value
		v.cas = atomic.AddUint32(&s.cas, 1)

		return api.WasmResultOk
	}

	s.m[key] = &sharedDataItem{
		data: value,
		cas:  atomic.AddUint32(&s.cas, 1),
	}

	return api.WasmResultOk
}

var globalSharedData = &sharedData{
	m: make(map[string]*sharedDataItem),
}

func (d *DefaultImportsHandler) GetSharedData(key string) (string, uint32, api.WasmResult) {
	return globalSharedData.get(key)
}

func (d *DefaultImportsHandler) SetSharedData(key string, value string, cas uint32) api.WasmResult {
	return globalSharedData.set(key, value, cas)
}

// shared queue

type sharedQueue struct {
	id    uint32
	name  string
	lock  sync.RWMutex
	queue []string
}

func (s *sharedQueue) enque(value string) api.WasmResult {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.queue = append(s.queue, value)

	return api.WasmResultOk
}

func (s *sharedQueue) deque() (string, api.WasmResult) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.queue) == 0 {
		return "", api.WasmResultEmpty
	}

	v := s.queue[0]
	s.queue = s.queue[1:]

	return v, api.WasmResultOk
}

type sharedQueueRegistry struct {
	lock             sync.RWMutex
	nameToIDMap      map[string]uint32
	m                map[uint32]*sharedQueue
	queueIDGenerator uint32
}

func (s *sharedQueueRegistry) register(queueName string) (uint32, api.WasmResult) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if queueID, ok := s.nameToIDMap[queueName]; ok {
		return queueID, api.WasmResultOk
	}

	newQueueID := atomic.AddUint32(&s.queueIDGenerator, 1)
	s.nameToIDMap[queueName] = newQueueID
	s.m[newQueueID] = &sharedQueue{
		id:    newQueueID,
		name:  queueName,
		queue: make([]string, 0),
	}

	return newQueueID, api.WasmResultOk
}

func (s *sharedQueueRegistry) resolve(queueName string) (uint32, api.WasmResult) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if queueID, ok := s.nameToIDMap[queueName]; ok {
		return queueID, api.WasmResultOk
	}

	return 0, api.WasmResultNotFound
}

func (s *sharedQueueRegistry) get(queueID uint32) *sharedQueue {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if queue, ok := s.m[queueID]; ok {
		return queue
	}

	return nil
}

var globalSharedQueueRegistry = &sharedQueueRegistry{
	nameToIDMap: make(map[string]uint32),
	m:           make(map[uint32]*sharedQueue),
}

func (d *DefaultImportsHandler) RegisterSharedQueue(queueName string) (uint32, api.WasmResult) {
	return globalSharedQueueRegistry.register(queueName)
}

func (d *DefaultImportsHandler) ResolveSharedQueue(queueName string) (uint32, api.WasmResult) {
	return globalSharedQueueRegistry.resolve(queueName)
}

func (d *DefaultImportsHandler) EnqueueSharedQueue(queueID uint32, data string) api.WasmResult {
	queue := globalSharedQueueRegistry.get(queueID)
	if queue == nil {
		return api.WasmResultNotFound
	}

	return queue.enque(data)
}

func (d *DefaultImportsHandler) DequeueSharedQueue(queueID uint32) (string, api.WasmResult) {
	queue := globalSharedQueueRegistry.get(queueID)
	if queue == nil {
		return "", api.WasmResultNotFound
	}

	return queue.deque()
}
