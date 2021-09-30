module shard_interface

go 1.16

require (
	src/config v0.0.0
	src/network v0.0.0
)

replace (
	src => ../../../src
	src/config => ../../config
	src/meta => ../../meta
	src/network => ../../network
)
