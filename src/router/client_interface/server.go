package client_interface

import (
	"crypto/cipher"
	"log"
	"net"

	"../settings"
)

type client_conn struct {
	net.Conn
	aes cipher.Block
}

func init_client() {
	PORT := ":" + settings.Client_port
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go Handle(c)
	}
}

func Send(conn client_conn) {

}
