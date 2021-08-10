package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func dial_server(router_addr string, mycfg []byte, Handler func(net.Conn, string), cfgfn func(map[string]interface{}, func(net.Conn, string)) error) error {

	defer wg.Done()

	conn, err := net.Dial("tcp", router_addr)

	if err != nil {
		log.Println(err)
		return err
	}

	defer conn.Close()

	send(conn, mycfg) // send my config to router, router reads and decides

	connbuff := bufio.NewReader(conn)

	config_json, err := recieve(connbuff) // read config sent from router

	if err != nil {
		log.Println(err)
		return err
	}

	var config map[string]interface{}
	json.Unmarshal([]byte(config_json), &config)

	err = cfgfn(config, Handler)

	if err != nil {
		log.Println(err)
		return err
	}

	if contains(current_router_ipv4, router_addr) { // make sure each router and shard only connected onece
		return errors.New("router has connected")
	}

	if v, ok := config["mac"]; ok {
		register_router(conn, v.(string), router_addr)

		defer closed_router(conn, v.(string), router_addr)

	} else if config["role"] == "shard" {
		shard_register_router(conn, router_addr)

		defer shard_closed_router(conn, router_addr)
	}

	for {
		message, err := recieve(connbuff)

		if err != nil {
			return err
		}

		go Handler(conn, message)
		wg.Add(1)
	}
}

func shardConfig(config map[string]interface{}, handler func(net.Conn, string)) error {

	if v, ok := config["router_ipv4"]; ok {

		var list []string
		for _, v := range v.([]interface{}) { // convert []interface{} to []string
			list = append(list, v.(string))
		}
		current_router_ipv4 = removeDuplicateStrings(list)

		for _, ip := range current_router_ipv4 {
			if ip != Settings.router_addr {
				go dial_server(ip, mycfg, ShardHandler, secondConfig)
				wg.Add(1)
			}

		}
	}
	return nil
}

func routerConfig(config map[string]interface{}, handler func(net.Conn, string)) error {
	if v, ok := config["router_ipv4"]; ok {

		var list []string
		for _, v := range v.([]interface{}) { // convert []interface{} to []string
			list = append(list, v.(string))
		}

		current_router_ipv4 = removeDuplicateStrings(list)
		for _, ip := range current_router_ipv4 {
			if ip != Settings.router_addr && ip != Settings.host+":"+Settings.port {
				go dial_server(ip, mycfg, RouterHandler, secondConfig)
				wg.Add(1)
			}

		}
	}
	return nil
}

func secondConfig(config map[string]interface{}, handler func(net.Conn, string)) error { // secondary configaration, connect to all routers
	if v, ok := config["router_ipv4"]; ok {
		var list []string
		for _, v := range v.([]interface{}) { // convert []interface{} to []string
			list = append(list, v.(string))
		}

		more_ip := difference(current_router_ipv4, list)

		more_ip = removeDuplicateStrings(more_ip)

		for _, ip := range more_ip {

			if ip != Settings.router_addr {
				go dial_server(ip, mycfg, handler, secondConfig)
				wg.Add(1)
			}

			current_router_ipv4 = append(current_router_ipv4, ip)
		}

	}
	return nil
}
