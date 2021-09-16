package client_interface

import (
	"crypto/cipher"
	"log"
	"net"
	"src/network"

	"../../settings"
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

func Send(conn client_conn, msg string) (int, error) {
	a, err := AESencrypt(conn.aes, msg)
	if err != nil {
		return 1, err
	}
	return network.Send(conn, []byte(a))
}

func Recv(conn client_conn) (string, error) {
	b, err := network.Recieve(conn)
	if err != nil {
		return "", err
	}
	c, err := AESdecrypt(conn.aes, string(b))
	if err != nil {
		return "", err
	}
	return c, nil
}
