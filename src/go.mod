module src

go 1.16

require (
	github.com/goccy/go-json v0.7.8
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/shirou/gopsutil v3.21.8+incompatible // indirect
	rogchap.com/v8go v0.6.0
	src/lazy v0.0.0
	src/network v0.0.0
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
	src/static => ./static
	src/tensorflow => ./tensorflow
	src/types => ./types
)
