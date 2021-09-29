package application_interface

import (
	"net"

	"../../network"
)

var JobQueue chan Job

var Jobs = map[uint64]Job{}

type Job struct {
	id     uint64
	client *net.Conn
	msg    []byte // Job does not directly store bytes
}

func Handle(conn net.Conn) {
	for {
		c, err := network.Recieve(conn)
		if err != nil {
			return
		}
		var d string
		if err != nil {
			d = err.Error()
		}
		network.Send(conn, []byte(d))
	}
}
