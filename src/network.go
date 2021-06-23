package router

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func dial_server(router_addr string, mycfg string, Handler func(net.Conn, string)) error {

	conn, err := net.Dial("tcp", router_addr)

	defer func() {
		conn.Close()
	}()

	if err != nil {
		return err
	}

	fmt.Fprintln(conn, mycfg) // send my config to router, router reads and decides

	connbuff := bufio.NewReader(conn)

	config_json, err := connbuff.ReadString('\n') // read config sent from router

	if err != nil {
		return err
	}

	var config map[string]interface{}
	json.Unmarshal([]byte(config_json), &config)

	for {
		message, err := connbuff.ReadString('\n')

		if err != nil {
			return err
		}

		log.Println(time.Now().String() + " Message from server: " + message)
		go Handler(conn, message)
	}
}
