package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

var command_query = map[string]string{}

func start_listening() error { // main loop

	ln, err := net.Listen("tcp", ":8080")
	fmt.Println("[server] server start")

	defer ln.Close()

	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		fmt.Println("connection from " + conn.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
		}
		go router_handleConnection(conn) // handle connection in new routine
	}
}

func router_handleConnection(conn net.Conn) { // this function handles a single connection

	connbuff := bufio.NewReader(conn)

	conf_json, err := connbuff.ReadString('\n') // read config from client

	CheckErr(err)

	var config map[string]interface{}

	err = json.Unmarshal([]byte(conf_json), &config) // decode client config
	CheckErr(err)

	role, mac, err := router_connection_config(conn, config) // handle client config and send config back to remote
	CheckErr(err)

	log.Println("[" + role + "]" + time.Now().String() + " connected")

	switch role {
	case "shard": // register shard

		register_shard(conn, mac)

		defer func() {
			closed_shard(conn, mac)
			conn.Close()
		}()

	case "router": // register router

		register_router(conn, mac)

		defer func() {
			closed_router(conn, mac)
			conn.Close()
		}()

	}

	for {
		message, err := connbuff.ReadString('\n')
		CheckErr(err)
		go RouterHandler(conn, message) // handle map sync and shard data feed back
	}
}

func router_connection_config(conn net.Conn, config map[string]interface{}) (string, string, error) {

	// read the configeration and determine the remote's role
	// validate password if neccesary
	var role string
	var mac string
	if role, ok := config["role"]; ok {
		switch role {

		case "shard":

			if pass, ok := config["pass"]; ok {
				if pass.(string) != password { // password is a global var, if not specified, default as ""
					fmt.Fprintln(conn, "Invalid password")
					conn.Close()
				}

			}
			role = "shard"
			mycfgmap := map[string]interface{}{"router_ipv4": current_router_ipv4}

			mycfg, _ := json.Marshal(mycfgmap)

			fmt.Fprint(conn, mycfg) // send config to remote

		case "router":

			if pass, ok := config["pass"]; ok {
				if pass.(string) != password { // password is a global var, if not specified, default as ""
					fmt.Fprintln(conn, "Invalid password")
					conn.Close()
				}

			}
			role = "router"

			mycfgmap := map[string]interface{}{"mac": MAC_Address,
				"router_load": router_load, "router_ipv4": current_router_ipv4}

			mycfg, _ := json.Marshal(mycfgmap)

			fmt.Fprint(conn, mycfg)

		case "client":

			role = "client"
			mycfgmap := map[string]interface{}{}

			mycfg, _ := json.Marshal(mycfgmap)

			fmt.Fprint(conn, mycfg)

		default:
			return "", "", errors.New("no role specified")

		}
	}
	return role, mac, nil
}
