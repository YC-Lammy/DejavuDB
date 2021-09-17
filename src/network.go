package main

import (
	"bufio"
	"errors"
	"log"
	"net"

	json "github.com/goccy/go-json"

	"src/lazy"
	"src/settings"
)

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

	if lazy.Contains(current_router_ipv4, router_addr) { // make sure each router and shard only connected onece
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
		current_router_ipv4 = lazy.RemoveDuplicateStrings(list)

		for _, ip := range current_router_ipv4 {
			if ip != settings.Router_addr {
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

		current_router_ipv4 = lazy.RemoveDuplicateStrings(list)
		for _, ip := range current_router_ipv4 {
			if ip != settings.Router_addr && ip != settings.Host+":"+settings.Port {
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

		more_ip := lazy.Difference_str_arr(current_router_ipv4, list)

		more_ip = lazy.RemoveDuplicateStrings(more_ip)

		for _, ip := range more_ip {

			if ip != settings.Router_addr {
				go dial_server(ip, mycfg, handler, secondConfig)
				wg.Add(1)
			}

			current_router_ipv4 = append(current_router_ipv4, ip)
		}

	}
	return nil
}
