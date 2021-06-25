package main

import (
	"fmt"
	"log"
	"net"
)

var shard_connected = 0

var router_connected = 0

var data_map = map[string]interface{}{}

var shard_map = map[string]net.Conn{} // mac addr : net.Conn or false if disconnected

var router_map = map[string]net.Conn{} //

var router_load = map[string]int{}

var current_router_ipv4 = []string{} //

var user_map = map[string]interface{}{ // user_map will not be exposed to the out front
	/*
		GID 1–99 are reserved for the system and application use.
		GID 100+ allocated for the user’s group.
		UIDs 1–99 are reserved for other predefined accounts.
		UID 100–999 are reserved by system for administrative and system accounts/groups.
		UID 1000–10000 are occupied by applications account.
		UID 10000+ are used for user accounts.
	*/
	"adm":     map[string]interface{}{"id": "1"},    // admin, nearest to root
	"sudo":    map[string]interface{}{"id": "27"},   // config permission, upgrade and maintainance
	"dev":     map[string]interface{}{"id": "30"},   // developers, view logs and cofigs
	"monitor": map[string]interface{}{"id": "80"},   // analystics, no admin permissions
	"user":    map[string]interface{}{"id": "100"},  // regular user, no additional permissions
	"public":  map[string]interface{}{"id": "1000"}} // public access, no authorization needed

func register_shard(conn net.Conn, mac string) {

	remote := conn.RemoteAddr().String()

	log.Println("[shard] " + remote + " Connected")
	shard_map[mac] = conn
	shard_connected += 1
}

func closed_shard(conn net.Conn, mac string) {

	remote := conn.RemoteAddr().String()

	log.Println("[shard] " + remote + " Disconnected")
	shard_map[mac] = nil

	shard_connected -= 1
}

func register_router(conn net.Conn, mac string, port string) {

	remote := conn.RemoteAddr().String()

	log.Println("[router] " + remote + " Connected")
	router_map[mac] = conn
	router_connected += 1

	current_router_ipv4 = append(current_router_ipv4, port)
}

func closed_router(conn net.Conn, mac string, port string) {

	remote := conn.RemoteAddr().String()

	log.Println("[router] " + remote + " Disconnected")
	router_map[mac] = nil

	router_connected -= 1

	current_router_ipv4 = removeItem(current_router_ipv4, port)
}

func send_to_all_router(message interface{}) {
	for _, v := range router_map {
		if v != nil {
			fmt.Fprintln(v, message)
		}

	}
}

func send_to_all_shard(message interface{}) {
	for _, v := range shard_map {
		if v != nil {
			fmt.Fprintln(v, message)
		}
	}
}
