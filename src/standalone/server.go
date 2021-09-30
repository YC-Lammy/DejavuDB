package standalone

import (
	"log"
	"net"
	"src/config"
	"src/standalone/client_interface"
)

func Start() {
	c, err := net.Listen("tcp", config.Client_port) // client interface
	if err != nil {
		panic(err)
	}
	defer c.Close()
	for {
		conn, err := c.Accept()
		if err != nil {
			log.Println(err)
		}
		go client_interface.Handle(conn)
	}
}
