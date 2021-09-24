module network

go 1.16

require (
	github.com/fxamacker/cbor/v2 v2.3.0
	github.com/goccy/go-json v0.7.8
	src/config v0.0.0
	src/meta v0.0.0
)

replace(
	src/config => ../config
	src/meta => ../meta
)