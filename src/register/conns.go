package register

import (
	"net"
)

type Conn_register struct{
	Conn net.Conn
	Buffer *bufio.Reader
	ID uint16
	Role string
}
var shard_connected = 0
var router_connected = 0

var Shards = map[uint16]*Conn_register{}
var Heartbeat_conns = []net.Conn{}
