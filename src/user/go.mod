module user

go 1.16

require (
	github.com/goccy/go-json v0.7.9
	rogchap.com/v8go v0.6.0 // indirect
	src/config v0.0.0
	src/lazy v0.0.0
)

replace (
	src/config => ../config
	src/lazy => ../lazy
)
