package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func dial_server(router_addr string, mycfg []byte, Handler func(net.Conn, string), cfgfn interface{}) error {

	conn, err := net.Dial("tcp", router_addr)

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

	switch v := cfgfn.(type) {
	case func(map[string]interface{}) error:
		err := v(config)
		if err != nil {
			log.Fatalln(err)
			return err
		}
	}

	if v, ok := config["mac"]; ok {
		register_router(conn, v.(string))
		defer func() {
			conn.Close()
			closed_router(conn, v.(string))
		}()

	}

	for {
		message, err := connbuff.ReadString('\n')

		if err != nil {
			return err
		}

		log.Println(time.Now().String() + " Message from server: " + message)
		go Handler(conn, message)
	}
}
