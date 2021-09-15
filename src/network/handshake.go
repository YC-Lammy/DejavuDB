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

func SendHandshake(conn net.Conn) error {
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
	if string(msg) == "ok" {
		return nil
	}
	return errors.New(string(msg))
}

func RecvHandshake(conn net.Conn) (Handshake, error) {
	msg, err := Recieve(conn)
	if err != nil {
		Send(conn, []byte(err.Error()))
		return Handshake{}, err
	}
	handshake := Handshake{}
	err = json.Unmarshal([]byte(msg), &handshake)
	if err != nil {
		Send(conn, []byte(err.Error()))
		return handshake, err
	}

	if handshake.Pass != settings.Password {
		Send(conn, []byte("password incorrect"))
		return handshake, errors.New("password incorrect")
	}

	register.Shards[handshake.ID] = &register.Conn_register{
		Conn: conn,
	}

	Send(conn, []byte("ok"))
	return handshake, nil
}
