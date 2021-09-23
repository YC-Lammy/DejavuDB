module javascriptAPI

go 1.16

require (
	github.com/goccy/go-json v0.7.8
	rogchap.com/v8go v0.6.0
	src v0.0.0
	src/types v0.0.0
)

replace(
	src => ../../src
	src/types => ../../src/types
)
