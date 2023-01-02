module github.com/banzaicloud/proxy-wasm-go-host/runtime/wazero

go 1.18

replace github.com/banzaicloud/proxy-wasm-go-host => ../../

require (
	github.com/banzaicloud/proxy-wasm-go-host v0.2.1-0.20221123073237-4f948bf02510
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.8.1
	github.com/tetratelabs/wabin v0.0.0-20220927005300-3b0fbf39a46a
	github.com/tetratelabs/wazero v1.0.0-pre.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
