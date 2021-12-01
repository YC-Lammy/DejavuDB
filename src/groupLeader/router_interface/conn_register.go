package router_interface

import (
	"net"

	"src/network"
)

var server_id uint64

var router_connected = 0

var Router_conns = map[uint64]Server_conn{}

type Server_conn struct {
	Conn      net.Conn
	Heartbeat net.Conn
	Meta      net.Conn
}

func Send_to_all_router(msg []byte) {
	for _, v := range Router_conns {
		network.Send(v.Conn, msg)
	}
}
