module shard_interface

go 1.16

require(
    src/network v0.0.0
)

replace(
    src/network => ../../network
)
