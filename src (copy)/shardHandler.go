package main

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/DmitriyVTitov/size"
)

func ShardHandler(conn net.Conn, message string) {
	defer wg.Done()

	processID := ""

	commands := strings.Split(message, " ")
	if commands[0] == "processID" {
		commands = commands[2:]
		processID = commands[0] + " " + commands[1]
	}
	var result []byte
	switch commands[0] {

	case "shardsize":
		result = []byte(strconv.FormatInt(int64(getShardSize()), 10))

	case "Set", "Update", "Delete", "Get", "Clone", "Move": // non-sql
		v, err := Nosql_Handler(commands)
		if err != nil {
			send(conn, []byte(processID+" "+err.Error()))
		} else {
			send(conn, []byte(processID+" "+*v))
		}

	case "SQL":

	case "connect":
		if !(contains(current_router_ipv4, commands[len(commands)-1])) {
			go dial_server(commands[len(commands)-1], mycfg, ShardHandler, secondConfig)
			wg.Add(1)
		}

	case "monitor":
		arr, err := json.Marshal(monitor())
		if err != nil {
			log.Println(err)
			return
		}
		send(conn, []byte("monitorResult "+string(arr)))

	default:
		send(conn, []byte(processID+" Invalid"))
		return
	}
	send(conn, result)
}

func getShardSize() int {
	return size.Of(shardData)
}
