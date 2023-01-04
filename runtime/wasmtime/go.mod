module github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmtime

go 1.18

replace github.com/banzaicloud/proxy-wasm-go-host => ../../

require (
	emperror.dev/errors v0.8.1
	github.com/banzaicloud/proxy-wasm-go-host v0.2.1-0.20221123073237-4f948bf02510
	github.com/bytecodealliance/wasmtime-go/v3 v3.0.2
	github.com/go-logr/logr v1.2.3
	k8s.io/klog/v2 v2.80.1
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)
