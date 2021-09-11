package sql

import "net"

func Listen() {
	tcp, err := net.Listen("tcp", ":1433")
	udp, err := net.Listen("udp", ":1434")
}
