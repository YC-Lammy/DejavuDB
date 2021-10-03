module standalone

go 1.16

require (
	src/config v0.0.0
	src/standalone/client_interface v0.0.0
)

replace (
	src/config => ../config
	src/datastore => ../datastore
	src/javascriptAPI => ../javascriptAPI
	src/yaegiAPI => ../../yaegiAPI
	src/lazy => ../lazy
	src/meta => ../meta
	src/network => ../network
	src/standalone/client_interface => ./client_interface
	src/types => ../types
	src/types/decimal => ../types/decimal
	src/types/float128 => ../types/float128
	src/types/int128 => ../types/int128
	src/user => ../user
)
