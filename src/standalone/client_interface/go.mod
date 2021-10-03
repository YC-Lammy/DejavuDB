module client_interface

go 1.16

require (
	github.com/goccy/go-json v0.7.9
	rogchap.com/v8go v0.6.0
	src/config v0.0.0
	src/javascriptAPI v0.0.0
	src/yaegiAPI v0.0.0
	src/lazy v0.0.0
	src/network v0.0.0-00010101000000-000000000000
	src/user v0.0.0
)

replace (
	src/config => ../../config
	src/datastore => ../../datastore
	src/javascriptAPI => ../../javascriptAPI
	src/yaegiAPI => ../../yaegiAPI
	src/lazy => ../../lazy
	src/meta => ../../meta
	src/network => ../../network
	src/standalone => ../../standalone
	src/standalone/client_interface => ./
	src/types => ../../types
	src/types/decimal => ../../types/decimal
	src/types/float128 => ../../types/float128
	src/types/int128 => ../../types/int128
	src/user => ../../user
)
