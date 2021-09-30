package client_interface

import (
	"crypto/aes"
	"net"

	"src/user"

	"src/network"

	json "github.com/goccy/go-json"

	"src/lazy"
	"src/router"
)

type user_ struct {
	Username string
	Password string
}

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
	u, err := Recv(con) // {Username:"name", Password:"password"}
	if err != nil {
		return
	}
	f := user_{}
	err = json.Unmarshal([]byte(u), &a)
	if err != nil {
		return
	}
	if u, ok := user.Login(f.Username, f.Password); ok {
		con.gid = u.Gid
		con.id = u.Id
	} else {
		return
	}

	for {
		c, err := Recv(con)
		if err != nil {
			return
		}
		router.NewJob(&con.Conn, c)
		if err != nil {
			c = []byte(err.Error())
		}
		_, err = Send(con, c)
		if err != nil {
			return
		}
	}
}
