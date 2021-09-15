package client_interface

import (
	"crypto/aes"
	"net"
	"src/network"

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
	conn = client_conn{Conn: conn, aes: aesk}
}
