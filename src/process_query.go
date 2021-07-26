package main

import (
	"net"
	"sync"
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

// adding more query maps speeds up the whole operation since the bottleneck is to add process
var process_query = map[int]*process{}

var process_query1 = map[int]*process{}

var process_query2 = map[int]*process{}

var process_query3 = map[int]*process{}

var process_query_lock = sync.Mutex{}

var process_query_lock1 = sync.Mutex{}

var process_query_lock2 = sync.Mutex{}

var process_query_lock3 = sync.Mutex{}

var process_id_counter int8 = 0

// a map can only be written once every time, this is to prevent concurrent map writting

func add_process(client net.Conn, responses int) int { // create process and return the id
	now := time.Now()
	id := get_process_id()
	newprocess := process{id: id, client: &client, timeout: now.Add(10 * time.Minute), result: []byte{}, responses: responses}
	switch process_id_counter {
	case 0:
		process_id_counter++
		process_query_lock.Lock()
		process_query[id] = &newprocess
		process_query_lock.Unlock()
	case 1:
		process_id_counter++
		process_query_lock1.Lock()
		process_query1[id] = &newprocess
		process_query_lock1.Unlock()
	case 2:
		process_id_counter++
		process_query_lock2.Lock()
		process_query2[id] = &newprocess
		process_query_lock2.Unlock()
	case 3:
		process_id_counter = 0
		process_query_lock3.Lock()
		process_query3[id] = &newprocess
		process_query_lock3.Unlock()
	}

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
				process_query_lock.Lock()
				delete(process_query, k)
				process_query_lock.Unlock()
			}
		}
	}
}
