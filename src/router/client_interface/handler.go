package client_interface

import (
	"crypto/aes"
	"net"
	"src/network"

	"../../javascriptAPI"
	"../../lazy"
)

func Handle(conn net.Conn) {
	defer conn.Close()
	k, err := network.Recieve(conn)
	if err != nil {
		return
	}
	key, err := ParseRsaPublicKeyFromPemStr(string(k))
	if err != nil {
		return
	}
	a := lazy.RandString(32)

	network.Send(conn, []byte(RSA_OAEP_Encrypt(a, *key)))
	aesk, err := aes.NewCipher([]byte(a))
	if err != nil {
		return
	}
	con := client_conn{Conn: conn, aes: aesk}

	for {
		c, err := Recv(con)
		if err != nil {
			return
		}
		c, err = javascriptAPI.Javascript_run_isolate(c)
		if err != nil {
			c = err.Error()
		}
		Send(con, c)
	}
}
