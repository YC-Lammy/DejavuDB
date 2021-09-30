module user

go 1.16

require (
    github.com/goccy/go-json v0.7.9
    src/lazy v0.0.0
    src/config v0.0.0
)

replace(
    src/lazy => ../lazy
    src/config => ../config
)
