package main

import (
	"log"
	"net"

	"src/lazy"
)

type shardDetail struct {
	size      int
	conn      net.Conn
	mem_load  int
	cpu_load  int
	disk_load int
	mem_size  int
	disk_size int
}

var shard_connected = 0

var router_connected = 0

var shard_map = map[string]net.Conn{} // mac addr : net.Conn or false if disconnected

var shard_size = map[string]int{}

var router_map = map[string]net.Conn{} //

var router_load = map[string]int{}

var current_router_ipv4 = []string{} //

var log_servers = []net.Conn{}

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

func shard_register_router(conn net.Conn, port string) {

	remote := conn.RemoteAddr().String()

	log.Println("[router] " + remote + " Connected")
	router_connected += 1

	current_router_ipv4 = append(current_router_ipv4, port)
}

func closed_router(conn net.Conn, mac string, port string) {

	remote := conn.RemoteAddr().String()

	log.Println("[router] " + remote + " Disconnected")
	router_map[mac] = nil

	router_connected -= 1

	current_router_ipv4 = lazy.RemoveItem(current_router_ipv4, port)
}

func shard_closed_router(conn net.Conn, port string) {

	remote := conn.RemoteAddr().String()

	log.Println("[router] " + remote + " Disconnected")

	router_connected -= 1

	current_router_ipv4 = lazy.RemoveItem(current_router_ipv4, port)
}

func register_log(conn net.Conn) {
	remote := conn.RemoteAddr().String()
	log_servers = append(log_servers, conn)
	log.Println("[log] " + remote + " Connected")
}

func closed_log(conn net.Conn) {

	remote := conn.RemoteAddr().String()

	for i, v := range log_servers {
		if v == conn {
			log_servers = append(log_servers[:i], log_servers[i+1:]...)
		}
	}
	log.Println("[log] " + remote + " Disconnected")
}

func send_to_all_router(message []byte) {
	for _, v := range router_map {
		if v != nil {
			send(v, message)
		}

	}
}

func send_to_all_shard(message []byte) {
	for _, v := range shard_map {
		if v != nil {
			send(v, message)
		}
	}
}
