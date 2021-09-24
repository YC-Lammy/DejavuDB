module javascriptAPI

go 1.16

require (
	github.com/goccy/go-json v0.7.8
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/traefik/yaegi v0.10.0
	go.mongodb.org/mongo-driver v1.7.2 // indirect
	rogchap.com/v8go v0.6.0
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
)
