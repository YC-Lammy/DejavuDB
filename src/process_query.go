package main

import (
	"net"
	"time"
)

type process struct {
	client net.Conn  // reply to client using this field
	time   time.Time // time when create process
}

var process_id int = 0 // process id is a router specific id, it does not represent the global process number

var process_query map[int]process

func add_process(client net.Conn) int {
	newprocess := process{client: client, time: time.Now()}
	id := get_process_id()
	process_query[id] = newprocess
	return id
}

func get_process_id() int {
	if process_id > 10000 { // the 0th process is most likely finished
		process_id = 0
	}
	process_id += 1
	return process_id
}
