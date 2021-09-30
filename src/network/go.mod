module network

go 1.16

require (
	github.com/fxamacker/cbor/v2 v2.3.0
	github.com/goccy/go-json v0.7.8
	rogchap.com/v8go v0.6.0 // indirect
	src/config v0.0.0
	src/meta v0.0.0
)

replace (
	src => ../../src
	src/config => ../config
	src/meta => ../meta
)
