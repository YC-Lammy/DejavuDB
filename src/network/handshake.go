package network

import (
	"errors"
	"net"

	json "github.com/goccy/go-json"

	"../register"
	"../settings"
)

type Handshake struct {
	Role string
	Pass string
	Host string
	Port string

	ID uint16 //0 if new born
}

func SendHandshake(conn *net.Conn) error {
	v := Handshake{
		Role: settings.Role,
		Pass: settings.Password,
		Host: settings.Host,
		Port: settings.Port,

		ID: settings.ID,
	}
	js, err := json.Marshal(v)
	if err != nil {
		return err
	}
	Send(conn, js)
	msg, err := Recieve(conn)
	if err != nil {
		return err
	}
	if msg == "ok" {
		return nil
	}
	return errors.New(msg)
}

func RecvHandshake(conn *net.Conn) error {
	msg, err := Recieve(conn)
	if err != nil {
		Send(conn, err.String())
		return err
	}
	handshake := Handshake{}
	err = json.Unmarshal([]byte(msg), &handshake)
	if err != nil {
		Send(conn, err.String())
		return err
	}

	if Pass != settings.Password {
		Send(conn, "password incorrect")
		return errors.New("password incorrect")
	}

	register.Shards[handshake.ID] = &register.Conn_register{
		Conn: conn,
	}

	Send(conn, "ok")
	return nil
}
