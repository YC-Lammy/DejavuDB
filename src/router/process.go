package router

import "net"

type process struct {
	id     uint64
	shard  *net.Conn
	client *net.Conn
}
