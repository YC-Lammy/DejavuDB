package standalone

import (
	"net"
	"src/config"
)

func Start() {
	c, err := net.Listen("tcp", config.Client_port) // client interface
	if err != nil {
		panic(err)
	}
}
