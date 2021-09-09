package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func Send(conn net.Conn, message []byte) (int, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint64(len(message)))
	if err != nil {
		return 0, err
	}
	fmt.Fprint(conn, string(buf.Bytes()))
	return fmt.Fprint(conn, string(message))
}

func Recieve(conn net.Conn) (string, error) {
	var length uint64
	var lenbuf = make([]byte, 8)
	conn.Read(lenbuf)
	buf := bytes.NewReader(lenbuf)
	binary.Read(buf, binary.LittleEndian, &length)
	message := make([]byte, length)
	_, err := conn.Read(message)
	if err != nil {
		return "", err
	}
	return string(message), nil
}
