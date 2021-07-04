package main

import (
	"net"
	"strings"
)

func RouterHandler(conn net.Conn, message string) { // this function handles any message recieved from client
	defer wg.Done()
	switch strings.Split(message, " ")[0] {
	case "client":
		message = strings.Replace(message, "client ", "", 1)
		router_clientHandler(conn, message)
	case "router":
	case "shard":
	}

}

func router_clientHandler(conn net.Conn, message string) {
	splited_message := strings.Split(message, " ")
	switch splited_message[0] {
	case "Set":
	case "Get":
	case "Update":
	case "Clone":
	case "Delete":
	case "Move":
	}
}
