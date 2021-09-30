module router

go 1.16

require (
	github.com/goccy/go-json v0.7.9 // indirect
	rogchap.com/v8go v0.6.0
	src/javascriptAPI v0.0.0
)

replace (
	src/config => ../config
	src/datastore => ../datastore
	src/javascriptAPI => ../javascriptAPI
	src/lazy => ../lazy
	src/meta => ../meta
	src/network => ../network
	src/types => ../types
	src/types/decimal => ../types/decimal
	src/types/float128 => ../types/float128
	src/types/int128 => ../types/int128
	src/user => ../user
)
