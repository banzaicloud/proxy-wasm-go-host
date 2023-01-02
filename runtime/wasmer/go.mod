module github.com/banzaicloud/proxy-wasm-go-host/runtime/wasmer

go 1.18

replace github.com/banzaicloud/proxy-wasm-go-host => ../..

require (
	github.com/banzaicloud/proxy-wasm-go-host v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.8.1
	github.com/wasmerio/wasmer-go v1.0.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)