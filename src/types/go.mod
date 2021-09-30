module types

go 1.16

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/traefik/yaegi v0.10.0 // indirect
	src/types/decimal v0.0.0
	src/types/float128 v0.0.0
	src/types/int128 v0.0.0
)

replace(
	src/types/decimal => ./decimal
	src/types/float128 => ./float128
	src/types/int128 => ./int128
)