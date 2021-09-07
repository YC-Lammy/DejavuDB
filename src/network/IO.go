package network

import(
	"net"

	"../register"
)

func Send(conn net.Conn, message []byte) (int, error) {
	fmt.Fprint(conn, int64(len(message)))
	return fmt.Fprint(conn, string(message))
}

func Recieve(buffer *bufio.Reader) (string, error) {
	var length [8]byte
	buffer.Read(length)
	message := make([]byte, int64(length))
	err := buffer.Read(message)
	if err != nil {
		return "", err
	}
	return string(message), nil
}