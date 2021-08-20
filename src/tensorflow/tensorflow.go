package tensorflow

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
)

var tf_model_register = map[string]*tfModel{}

var tf_python_server_conn *net.Conn

var tf_python_server_conn_buf *bufio.Reader

var Conn *net.Conn

var Connbuf *bufio.Reader

var Version string

type tfModel struct {
	name         string
	layer_count  int //keras
	param_count  int //keras
	input_shape  []int
	output_shape []int
	constructer  string //javascript

	path string
	lock sync.Mutex
}

func Init_tensorflow() error {
	if err := init_python_server(); err != nil {
		return err
	}
	co, err := net.Dial("localhost", "3247")
	if err != nil {
		return err
	}
	Conn = &co
	Connbuf = bufio.NewReader(co)

	return nil
}

func Tf_send(msg []byte) {
	header := strconv.Itoa(len(msg))
	for len(header) < 64 { // header must be length of 64
		header = "0" + header
	}
	fmt.Fprint(*Conn, header)
	fmt.Fprint(*Conn, msg)
}

func Tf_recv() ([]byte, error) {
	header := []byte{}
	for i := 0; i < 64; i++ {
		by, err := Connbuf.ReadByte()
		if err != nil { //error when no byte is avaliable
			i--      // do not count the loop
			continue // skip to next loop
		}
		header = append(header, by)
	}
	msg_l, err := strconv.Atoi(string(header))
	if err != nil {
		return nil, err
	}
	msg := []byte{}
	for i := 0; i < msg_l; i++ {
		by, err := Connbuf.ReadByte()
		if err != nil { //error when no byte is avaliable
			i--      // do not count the loop
			continue // skip to next loop
		}
		msg = append(msg, by)
	}
	return msg, nil
}
func tf_model_predict(model_name string, data interface{}) string {
	return ``
}

func Get_model_by_name(name string) (*tfModel, error) {
	if v, ok := tf_model_register[name]; ok {
		return v, nil
	}
	return nil, errors.New("model " + name + " does not exist")
}
