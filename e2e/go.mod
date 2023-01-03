module github.com/banzaicloud/proxy-wasm-go-host/e2e

go 1.18

replace (
	github.com/banzaicloud/proxy-wasm-go-host => ../
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmer => ../runtime/wasmer
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wazero => ../runtime/wazero
)

require (
	github.com/banzaicloud/proxy-wasm-go-host v0.2.1-0.20221123073237-4f948bf02510
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmer v0.0.0-00010101000000-000000000000
	github.com/banzaicloud/proxy-wasm-go-host/runtime/wazero v0.0.0-00010101000000-000000000000
	github.com/tetratelabs/wabin v0.0.0-20220927005300-3b0fbf39a46a
	mosn.io/pkg v1.3.0
)

require (
	emperror.dev/errors v0.8.1 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/tetratelabs/wazero v1.0.0-pre.4 // indirect
	github.com/wasmerio/wasmer-go v1.0.4 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	mosn.io/api v1.3.0 // indirect
)
