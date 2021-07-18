package main

import (
	"net"
	"time"
)

type process struct {
	id        int
	client    *net.Conn // reply to client using this field
	time      time.Time // time when create process
	responses int       // expected number of shard to be responding
	timeout   time.Time
	result    []byte
}

var process_id int = 0 // process id is a router specific id, it does not represent the global process number

var process_query = map[int]*process{}

var process_sync = make(chan *process)

// a map can only be written once every time, this is to prevent concurrent map writting
func process_query_sync() {
	for {
		newprocess := <-process_sync
		process_query[newprocess.id] = newprocess
	}
}

func add_process(client net.Conn, responses int) int { // create process and return the id
	now := time.Now()
	id := get_process_id()
	newprocess := process{id: id, client: &client, timeout: now.Add(10 * time.Minute), result: []byte{}, responses: responses}
	process_sync <- &newprocess
	return id
}

func get_process_id() int {
	if process_id > 10000 { // the 0th process is most likely finished
		process_id = 0
	}
	process_id += 1
	return process_id
}

func process_timeout_checker() {
	for {
		now := time.Now()
		time.Sleep(1 * time.Second)
		for k, v := range process_query {
			if v.timeout.Sub(now) < 0 {
				send(*v.client, []byte("process timeout"))
				delete(process_query, k)
			}
		}
	}
}
