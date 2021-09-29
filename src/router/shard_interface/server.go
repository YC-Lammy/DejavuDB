package shard_interface

import (
	"net"
	"src/network"
)

type Server_conn struct {
	Conn      net.Conn
	Heartbeat net.Conn
	Meta      net.Conn
}

var Shard_connected uint64
var Shard_conns = map[uint64]Server_conn{}

func Send_to_all_shard(msg []byte) {
	for _, v := range Shard_conns {
		network.Send(v.Conn, msg)
	}
}
