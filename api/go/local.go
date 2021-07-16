package dejavuDB

import (
	"bufio"
	"fmt"
	"net"
)

func send(conn net.Conn, message []byte) (int, error) {
	message = append(message, 0x00) // nul to mark end of section
	return fmt.Fprint(conn, string(message))
}

func recieve(buffer *bufio.Reader) (string, error) {
	message, err := buffer.ReadBytes(0x00)
	if err != nil {
		return "", err
	}
	return string(message[:len(message)-1]), nil
}

func contains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}
