package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
)

func RouterHandler(conn net.Conn, message string) { // this function handles any message recieved from client

	defer wg.Done()

	switch strings.Split(message, " ")[0] {

	case "CLIENT":
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
					send(v, []byte(id_str+" "+message)) // send to shard
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

func getShardMac(location string) ([]string, error) { // get the mac addr of the shard that saves the data

	keys := strings.Split(location, ".")
	var pointer map[string]interface{}

	if v, ok := data_map[keys[0]]; ok {
		if i, ok := v.([]string); ok {
			return i, nil
		}
		if i, ok := v.(map[string]interface{}); ok {
			pointer = i
		} else {
			return nil, errors.New("type not match")
		}
	}

	for _, key := range keys[1:] {
		if v, ok := pointer[key]; ok {
			switch v := v.(type) {
			case map[string]interface{}:
				pointer = v

			case []string:
				return v, nil

			}

		}
	}
	buffer := []map[string]interface{}{}
	macs := []string{}
	// mac not found, find every mac under the pointer instead
	for _, v := range pointer {
		switch v := v.(type) {
		case []string:
			macs = append(macs, v...)
		case map[string]interface{}:
			buffer = append(buffer, v)
		default:
			return nil, errors.New("invalid type")

		}
	}
	if len(buffer) > 0 {
		for len(buffer) > 0 {
			for i, v := range buffer {
				buffer = append(buffer[:i], buffer[i+1:]...)
				for _, v := range v {
					switch v := v.(type) {
					case []string:
						macs = append(macs, v...)
					case map[string]interface{}:
						buffer = append(buffer, v)
					default:
						return nil, errors.New("invalid type")
					}
				}
			}
		}
	}
	return macs, nil
}
