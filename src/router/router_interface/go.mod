module router_interface

go 1.17

require(
    src/network v0.0.0
)

replace(
    src/network => ../../network
)
