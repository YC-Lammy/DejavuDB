package clustermanager

import "net"

var server_id uint64

var shard_connected = 0
var router_connected = 0

var Shard_conns = map[uint64]Server_conn{}
var Router_conns = map[uint64]Server_conn{}

type Server_conn struct {
	Conn      net.Conn
	Heartbeat net.Conn
}
