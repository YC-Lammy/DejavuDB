package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"
)

func RouterHandler(conn net.Conn, message string) { // this function handles any message recieved from client

	splited := strings.Split(message, " ")

	defer wg.Done()

	switch splited[0] {

	case "CLIENT":
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

	case "processID": // shard returns value

		id, err := strconv.ParseInt(strings.Split(message, " ")[1], 10, 64)

		if err != nil {
			log.Println(err)
			return
		}
		send(process_query[int(id)].client, []byte(strings.Join(strings.Split(message, " ")[2:], " ")))

	case "monitor":
		arr, err := json.Marshal(monitor())
		if err != nil {
			log.Println(err)
			return
		}
		send(conn, []byte("monitorResult "+string(arr)))

	case "monitorResult": // recieve moniter values

		register_monitor_value(message[14:]) // monitor.go

	default: // send to shard
		router_apiHandler(conn, message) // handle request to shard
	}

}

func router_apiHandler(conn net.Conn, message string) {
	splited_message := strings.Split(message, " ")
	switch splited_message[0] {
	case "Get", "Update", "Delete":

		id := add_process(conn) // register process

		macs, err := getShardMac(splited_message[1])
		if err != nil {
			delete(process_query, id)
			log.Println(err)
			return
		}
		id_str := strconv.FormatInt(int64(id), 10)
		counter := 0

		for _, v := range macs {
			if v, ok := shard_map[v]; ok { // get net.Conn
				if v != nil {
					send(v, []byte("processID "+id_str+" "+message)) // send to shard
					counter += 1
				}
			}
		}
		if counter == 0 {
			delete(process_query, id)
		}

	case "Set":

	case "Clone", "Move":

	case "SQL":

	}
}

func router_clientHandler(conn net.Conn, message string) {

	result, err := execute_command(conn, message) // console command.go

	if err != nil {
		send(conn, []byte(err.Error()))
		return
	}
	send(conn, []byte(result))
}
