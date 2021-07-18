package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
)

var command_query = map[string]string{}

var invalid_password = map[string]int{} // wrong password mac, attemp times

func start_listening() error { // main loop
	defer wg.Done()

	ln, err := net.Listen("tcp", hostport)
	fmt.Println("[server] server start")

	//defer ln.Close()

	if err != nil {
		return err
	}

	c := color.New(color.FgHiRed).Add(color.Bold)
	c.Println("\nListening at " + hostport + "\n")

	for {
		conn, err := ln.Accept()
		//fmt.Println("connection from " + conn.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
		}
		go router_handleConnection(conn) // handle connection in new routine
		wg.Add(1)

	}
}

func router_handleConnection(conn net.Conn) { // this function handles a single connection

	defer conn.Close()
	defer wg.Done()

	connbuff := bufio.NewReader(conn)

	conf_json, err := recieve(connbuff) // read config from client

	if err != nil {
		log.Println(err)
		return
	}

	var config map[string]interface{}

	err = json.Unmarshal([]byte(conf_json), &config) // decode client config
	if err != nil {
		log.Println(err.(*json.SyntaxError))
		return
	}

	role, mac, port, err := router_connection_config(conn, config) // handle client config and send config back to remote
	if err != nil {
		log.Println(err)
		return
	}

	switch role {
	case "shard": // register shard

		register_shard(conn, mac)

		defer closed_shard(conn, mac)

	case "router": // register router

		register_router(conn, mac, port)

		defer closed_router(conn, mac, port)

		send_to_all_shard([]byte("connect " + port)) // call all shards to connect to new router

	case "client":

	case "log":
		register_log(conn)
		defer closed_log(conn)

	}

	for {
		message, err := recieve(connbuff)
		if err != nil {
			return
		}
		go RouterHandler(conn, message) // handle map sync and shard data feed back
		wg.Add(1)
	}
}

func router_connection_config(conn net.Conn, config map[string]interface{}) (string, string, string, error) {

	// read the configeration and determine the remote's role
	// validate password if neccesary
	var mac string
	if role, ok := config["role"]; ok {
		switch role {

		case "shard":

			mac := config["mac"].(string)
			port := config["port"].(string)

			if pass, ok := config["pass"]; ok {
				if pass.(string) != password { // password is a global var, if not specified, default as ""
					send(conn, []byte("Invalid password"))
					if _, ok := invalid_password[mac]; ok {
						invalid_password[mac] += 1
					} else {
						invalid_password[mac] = 1
					}
					conn.Close()
					return "", "", "", errors.New("Invalid password from " + mac)
				}

			}
			rol := "shard"

			mycfgmap := map[string]interface{}{"router_ipv4": current_router_ipv4}

			mycfg, _ := json.Marshal(mycfgmap)

			send(conn, mycfg) // send config to remote
			return rol, mac, port, nil

		case "router":

			mac := config["mac"].(string)
			port := config["port"].(string)

			if pass, ok := config["pass"]; ok {
				if pass.(string) != password { // password is a global var, if not specified, default as ""
					send(conn, []byte("Invalid password"))
					conn.Close()
					return "", "", "", errors.New("Invalid password from " + mac)
				}

			}
			rol := "router"

			mycfgmap := map[string]interface{}{"mac": MAC_Address,
				"router_load": router_load, "router_ipv4": current_router_ipv4}

			mycfg, _ := json.Marshal(mycfgmap)

			send(conn, mycfg) // send mycfg to remote
			return rol, mac, port, nil

		case "client":

			rol := "client"
			mycfgmap := map[string]interface{}{}

			mycfg, _ := json.Marshal(mycfgmap)

			send(conn, mycfg)

			return rol, "", "", nil

		case "log":

			port := config["port"].(string)

			if pass, ok := config["pass"]; ok {
				if pass.(string) != password { // password is a global var, if not specified, default as ""
					send(conn, []byte("Invalid password"))
					if _, ok := invalid_password[mac]; ok {
						invalid_password[mac] += 1
					} else {
						invalid_password[mac] = 1
					}
					conn.Close()
					return "", "", "", errors.New("invalid password")
				}

			}
			rol := "log"

			mycfgmap := map[string]interface{}{"router_ipv4": current_router_ipv4}

			mycfg, _ := json.Marshal(mycfgmap)

			send(conn, mycfg) // send config to remote
			return rol, "", port, nil

		default:
			return "", "", "", errors.New(conn.RemoteAddr().String() + " no role specified")

		}
	}
	return "", "", mac, nil
}
