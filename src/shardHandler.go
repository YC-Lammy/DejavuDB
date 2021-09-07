package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"src/lazy"
	"strconv"
	"strings"

	"src/settings"

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
	processes := strings.Split(strings.Join(commands, " "), "&&")
	var result = []byte{}
	for i, v := range processes {

		if settings.Debug {
			fmt.Println(v)
		}
		commands = strings.Split(v, " ")
		if len(commands) == 0 {
			continue
		}

		for commands[0] == "" { // remove all empty string
			commands = commands[1:]
			if len(commands) == 0 {
				continue
			}
		}

		switch commands[0] {

		case "shardsize":
			result = append(result, []byte(strconv.Itoa(getShardSize()))...)

		case "SQL":

			message = message[4:] // remove "SQL "

			if strings.Contains(strings.ToUpper(message), "SELECT") { // if it contains select statement, a table will be return

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
				json, err := table.Json()
				if err != nil {
					send(conn, []byte(processID+" "+err.Error()))
					return
				}
				if i != 0 {
					result = append(result, []byte(json+"\n")...)
				} else {
					result = append(result, []byte(processID+" "+json+"\n")...)
				}

			} else {

				r, err := sqliteDB.Exec(strings.Join(commands[1:], " "))
				if err != nil {
					send(conn, []byte(processID+" "+err.Error()))
					return
				}
				a, err := r.RowsAffected()
				if err != nil {
					send(conn, []byte(processID+" "+err.Error()))
					return
				}
				var b string

				if settings.Debug {
					c, err := r.LastInsertId()
					if err != nil {
						b = ""
					} else {
						b = " last insert Id:" + strconv.Itoa(int(c))
					}
				}
				if i != 0 {
					result = append(result, []byte(" rows affected:"+strconv.Itoa(int(a))+b+"\n")...)
				} else {
					result = append(result, []byte(processID+" rows affected:"+strconv.Itoa(int(a))+b+"\n")...)
				}

			}

		case "connect":
			if !(lazy.Contains(current_router_ipv4, commands[len(commands)-1])) {
				go dial_server(commands[len(commands)-1], mycfg, ShardHandler, secondConfig)
				wg.Add(1)
			}

		case "monitor":
			arr, err := json.Marshal(monitor())

			result = []byte("monitorResult " + string(arr))
			if err != nil {
				log.Println(err)
				result = []byte("monitorResult {}")
			}

		case "":
			continue

		default: // nosql
			v, err := Nosql_Handler(commands) // nosql handler will handlers all the execution
			if err != nil {
				send(conn, []byte(processID+" "+err.Error()))
				return
			}

			result = append(result, []byte(processID+" "+*v+"\n")...)
			processID = "" // purge the process id

		}
		if i > 40 {
			break
		}
	}
	send(conn, result)
}

func getShardSize() int {
	return size.Of(shardData)
}
