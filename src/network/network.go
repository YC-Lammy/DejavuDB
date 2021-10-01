package network

import (
	"net"
	"unsafe"
)

func Send(conn net.Conn, message []byte) (int, error) {
	l := uint64(len(message))

	conn.Write((*(*[8]byte)(unsafe.Pointer(&l)))[:])
	return conn.Write(message)
}

func Recieve(conn net.Conn) ([]byte, error) {
	var length uint64
	var lenbuf = make([]byte, 8)
	conn.Read(lenbuf)
	leng := [8]byte{}
	copy(leng[:], lenbuf)
	length = *(*uint64)(unsafe.Pointer(&leng))
	message := make([]byte, length)
	_, err := conn.Read(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
