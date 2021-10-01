package client_interface

import (
	"crypto/aes"
	"fmt"
	"net"
	"strings"

	"src/config"
	"src/user"

	"src/network"

	"src/lazy"
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
	con := Client_conn{Conn: conn, aes: aesk}
	u, err := Recv(con) // {Username:"name", Password:"password"}
	if err != nil {
		fmt.Println(err)
		return
	}
	up := strings.Split(string(u), ",")
	if config.Debug {
		fmt.Println(up)
	}
	if u, ok := user.Login(up[0], up[1]); ok {
		con.gid = u.Gid
		con.id = u.Id
		Send(con, []byte("login sucess"))
	} else {
		Send(con, []byte("login failed"))
		return
	}

	for {
		c, err := Recv(con)
		if err != nil {
			return
		}
		if config.Debug {
			fmt.Println("NewJob : " + string(c))
		}
		NewJob(&con, c)
	}
}
