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
	//fmt.Println(message)
	if commands[0] == "processID" {
		processID = commands[0] + " " + commands[1]
		commands = commands[2:]
	}
	var result []byte
	switch commands[0] {

	case "shardsize":
		result = []byte(strconv.Itoa(getShardSize()))

	case "Set", "Update", "Delete", "Get", "Clone", "Move", "Sizeof", "SizeOf", "Typeof", "TypeOf": // non-sql
		v, err := Nosql_Handler(commands)
		if err != nil {
			send(conn, []byte(processID+" "+err.Error()))
		} else {
			send(conn, []byte(processID+" "+*v))
		}

		return

	case "SQL":

		message = message[4:] // remove "SQL "
		result := ""

		if strings.Contains(strings.ToUpper(message), "SELECT ") {

			rows, err := sqliteDB.Query(strings.Join(commands[1:], " "))
			if err != nil {
				send(conn, []byte(processID+" "+err.Error()))
				return
			}

			if rows == nil {
				send(conn, []byte(processID+" variable name does ot exist"))
				return
			}

			table, err := read_SQL_Rows(rows)
			rows.Close()
			if err != nil {
				send(conn, []byte(processID+" "+err.Error()))
				return
			}
			result, err = table.Json()
			if err != nil {
				send(conn, []byte(processID+" "+err.Error()))
				return
			}

		} else {

			_, err := sqliteDB.Exec(strings.Join(commands[1:], " "))
			if err != nil {
				send(conn, []byte(processID+" "+err.Error()))
				return
			}
		}

		send(conn, []byte(processID+" "+result))
		return

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
