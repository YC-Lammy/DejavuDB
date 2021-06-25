package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func dial_server(router_addr string, mycfg []byte, Handler func(net.Conn, string), cfgfn func(map[string]interface{}) error) error {

	defer wg.Done()

	conn, err := net.Dial("tcp", router_addr)

	if err != nil {
		log.Println(err)
		return err
	}

	defer conn.Close()

	time.Sleep(10 * time.Millisecond)

	fmt.Fprintln(conn, string(mycfg)) // send my config to router, router reads and decides

	connbuff := bufio.NewReader(conn)

	config_json, err := connbuff.ReadString('\n') // read config sent from router

	if err != nil {
		log.Println(err)
		return err
	}

	var config map[string]interface{}
	json.Unmarshal([]byte(config_json), &config)

	err = cfgfn(config)

	if err != nil {
		log.Println(err)
		return err
	}

	if v, ok := config["mac"]; ok {
		register_router(conn, v.(string), router_addr)

		defer closed_router(conn, v.(string), router_addr)

	}

	for {
		message, err := connbuff.ReadString('\n')

		if err != nil {
			return err
		}

		log.Println("Message from server: " + message)

		go Handler(conn, message)
		wg.Add(1)
	}
}
