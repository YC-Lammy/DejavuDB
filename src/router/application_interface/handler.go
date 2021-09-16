package application_interface

import (
	"net"

	"../../javascriptAPI"
	"../../network"
)

func Handle(conn net.Conn) {
	for {
		c, err := network.Recieve(conn)
		if err != nil {
			return
		}
		var d string
		d, err = javascriptAPI.Javascript_run_isolate(string(c))
		if err != nil {
			d = err.Error()
		}
		network.Send(conn, []byte(d))
	}
}
