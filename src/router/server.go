package router

import (
	"net"
	"src/config"
)

func Start() {
	c, err := net.Listen("tcp", config.Client_port) // client interface
	if err != nil {
		panic(err)
	}
	r, err := net.Listen("tcp", ":5734") // router interface
	if err != nil {
		panic(err)
	}
	s, err := net.Listen("tcp", ":18241") // shard interface
	if err != nil {
		panic(err)
	}
	m, err := net.Listen("tcp", ":5734") // meta inteface
	if err != nil {
		panic(err)
	}
	if config.Leader_addr != "" {
		net.Dial("tcp", config.Leader_addr+":5734")
	}
}
