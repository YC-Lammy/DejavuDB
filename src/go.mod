module github.com/YC-Lammy/DejavuDB

go 1.16

require (
	rogchap.com/v8go v0.6.0
	src/config v0.0.0
	src/javascriptAPI v0.0.0
	src/lazy v0.0.0
	src/standalone v0.0.0
	src/static v0.0.0
)

replace (
	src/config => ./config
	src/contract => ./types/type-contract
	src/datastore => ./datastore
	src/javascriptAPI => ./javascriptAPI
	src/lazy => ./lazy
	src/meta => ./meta
	src/network => ./network
	src/sql => ./sql
	src/standalone => ./standalone
	src/standalone/client_interface => ./standalone/client_interface
	src/static => ./static
	src/tensorflow => ./tensorflow
	src/types => ./types
	src/types/decimal => ./types/decimal
	src/types/float128 => ./types/float128
	src/types/int128 => ./types/int128
	src/user => ./user
	src/yaegiAPI => ./yaegiAPI
)
