package dejavuDB

import (
	"bufio"
	"errors"
	"net"
	"strings"
)

var list_of_commands = []string{"Set", "Get", "Delete", "Update", "Clone", "Move",
	"useradd", "groupadd", "Login", "Logout",
	"atop", "cat", "cp", "chmod", "chown", "chgrp", "df", "dstat", "find",
	"free", "id", "last", "mv", "netstat", "rm", "sort", "w", "top", "tar"}

type admin struct {
	username   string
	token      string
	connection net.Conn
	connbuf    *bufio.Reader
}

func (ad *admin) Send(message string) (int, error) {
	return send(ad.connection, []byte("CLIENT "+message))
}

func (ad *admin) Login(usr, password string) error {
	_, err := send(ad.connection, []byte("CLIENT Login "+usr+" "+password))
	if err != nil {
		return err
	}
	r, err := recieve(ad.connbuf)
	if err != nil {
		return err
	}
	if r == "invalid" {
		return errors.New("incorrect username or password")
	}
	ad.token = r

	return nil
}

func (ad *admin) Useradd(args ...string) error {
	arg := ""
	if len(args) > 0 {
		arg = " " + strings.Join(args, " ")
	}
	_, err := ad.Send("useradd" + arg)
	if err != nil {
		return err
	}
	r, err := recieve(ad.connbuf)
	if err != nil {
		return err
	}
	if r != "" {
		return errors.New(r)
	}
	return nil
}

func (ad *admin) Groupadd(args ...string) error {
	arg := ""
	if len(args) > 0 {
		arg = " " + strings.Join(args, " ")
	}
	_, err := ad.Send("groupadd" + arg)
	if err != nil {
		return err
	}
	r, err := recieve(ad.connbuf)
	if err != nil {
		return err
	}
	if r != "" {
		return errors.New(r)
	}
	return nil
}
