module github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmtime/v3

go 1.18

replace github.com/banzaicloud/proxy-wasm-go-host => ../../

require (
	emperror.dev/errors v0.8.1
	github.com/banzaicloud/proxy-wasm-go-host v1.0.1
	github.com/bytecodealliance/wasmtime-go/v3 v3.0.2
	github.com/go-logr/logr v1.2.3
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.8.1
	github.com/tetratelabs/wabin v0.0.0-20220927005300-3b0fbf39a46a
	k8s.io/klog/v2 v2.90.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
