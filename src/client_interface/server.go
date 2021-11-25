package client_interface

import (
	"crypto/cipher"
	"log"
	"net"

	"dejavuDB/src/config"
	"dejavuDB/src/network"
)

type Client_conn struct {
	net.Conn
	aes cipher.Block
	id  uint32
	gid uint32
}

func init_client() {
	PORT := ":" + config.Client_port
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

func Send(conn Client_conn, msg []byte) (int, error) {
	a, err := AESencrypt(conn.aes, msg)
	if err != nil {
		return 1, err
	}
	return network.Send(conn, a)
}

func Recv(conn Client_conn) ([]byte, error) {
	b, err := network.Recieve(conn)
	if err != nil {
		return nil, err
	}
	c, err := AESdecrypt(conn.aes, b)
	if err != nil {
		return nil, err
	}
	return c, nil
}