module github.com/banzaicloud/proxy-wasm-go-host/e2e

go 1.18

replace (
	github.com/banzaicloud/proxy-wasm-go-host => ../
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmer => ../runtime/wasmer
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmtime/v3 => ../runtime/wasmtime
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wazero => ../runtime/wazero
)

require (
	github.com/banzaicloud/proxy-wasm-go-host v1.0.1
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmer v1.0.4-c0
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmtime/v3 v3.0.2-c0
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wazero v1.0.0-pre.8-c0
	github.com/tetratelabs/wabin v0.0.0-20220927005300-3b0fbf39a46a
)

require (
	emperror.dev/errors v0.8.1 // indirect
	github.com/bytecodealliance/wasmtime-go/v3 v3.0.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tetratelabs/wazero v1.2.1 // indirect
	github.com/wasmerio/wasmer-go v1.0.4 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	k8s.io/klog/v2 v2.90.0 // indirect
)
