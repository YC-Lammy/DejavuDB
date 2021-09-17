package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
	"net"
	"strconv"
	"strings"

	json "github.com/goccy/go-json"
)

var default_SQL_server = 0

func RouterHandler(conn net.Conn, message string) { // this function handles any message recieved from client

	splited := strings.Split(message, " ")

	defer wg.Done()

	switch splited[0] {

	case "CLIENT": // admin api command
		if splited[1] == "RSA" {
			priv, _ := rsa.GenerateKey(rand.Reader, 128) // skipped error checking for brevity
			pub, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
			if err != nil {
				send(conn, []byte(err.Error()))
			}
			send(conn, pub)
			return
		}
		router_clientHandler(conn, message[7:]) // handle system console users

	case "processID": // shard returns result

		id, err := strconv.ParseInt(splited[1], 10, 64)

		if err != nil {
			log.Println(err)
			return
		}
		for i, queries := range []*map[int]*process{&process_query, &process_query1, &process_query2, &process_query3} {
			query := *queries
			switch i {
			case 0:
				defer process_query_lock.Unlock()
				process_query_lock.Lock()
			case 1:
				defer process_query_lock1.Unlock()
				process_query_lock1.Lock()
			case 2:
				defer process_query_lock2.Unlock()
				process_query_lock2.Lock()
			case 3:
				defer process_query_lock3.Unlock()
				process_query_lock3.Lock()
			}
			if p, ok := query[int(id)]; ok {
				p.result = append(p.result, []byte(strings.Join(strings.Split(message, " ")[2:], " ")+"\n")...)

				p.responses = p.responses - 1

				if p.responses == 0 { // all shard responsed
					send(*p.client, p.result)
					delete(*queries, int(id))
				}
				return
			}
		}

	case "monitor": // request to send monitor values
		arr, err := json.Marshal(monitor())
		if err != nil {
			log.Println(err)
			return
		}
		send(conn, []byte("monitorResult "+string(arr)))

	case "monitorResult": // recieve moniter values

		register_monitor_value(message[14:]) // monitor.go

	case "registerKey":

	default: // send to shard, common api handler
		router_apiHandler(conn, message) // handle request to shard
	}

}

func router_apiHandler(conn net.Conn, message string) {

	if shard_connected == 1 { // skip all mapping process if only one shard
		id := add_process(conn, 1)
		id_str := strconv.Itoa(id)
		if message == "" {
			return
		}
		send_to_all_shard([]byte("processID " + id_str + " " + message))
		return
	}

	splited_message := strings.Split(message, " ")

	switch splited_message[0] {

	case "Get", "Update", "Delete": // only one address needed to be checked, no interaction involved
		conns, err := getShardConn(splited_message[1])
		id := add_process(conn, len(conns)) // register process

		if err != nil {
			delete(process_query, id)
			log.Println(err)
			return
		}
		id_str := strconv.FormatInt(int64(id), 10)
		counter := 0

		for _, v := range conns {
			_, err := send(v, []byte("processID "+id_str+" "+message)) // send to shard
			if err != nil {
				continue
			}
			counter += 1
		}
		if counter == 0 {
			delete(process_query, id)
		}

	case "Set": // one address to be checked if key has been used or not

	case "Clone", "Move":

	case "SQL":
		id := add_process(conn, 1)
		id_str := strconv.Itoa(id)
		for _, v := range data_replicates_groups[default_SQL_server].connections {
			send(v, []byte("processID "+id_str+" "+message))
		}

	}
}

func router_clientHandler(conn net.Conn, message string) { // execute logged in commands

	result, err := execute_command(conn, message) // console command.go

	if err != nil {
		if err.Error() == "do not send" {
			return
		}
		send(conn, []byte(err.Error()))
		return
	}
	send(conn, []byte(result))
}
