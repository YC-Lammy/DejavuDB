module yaegiAPI

go 1.16

require (
	github.com/traefik/yaegi v0.10.0
	src/datastore v0.0.0
	src/types v0.0.0
)

replace (
	src => ../../src
	src/config => ../config
	src/datastore => ../datastore
	src/lazy => ../../src/lazy
	src/network => ../network
	src/static => ../static
	src/types => ../types
	src/types/decimal => ../types/decimal
	src/types/float128 => ../types/float128
	src/types/int128 => ../types/int128
)
