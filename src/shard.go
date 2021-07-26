package main

import (
	"sync"
)

type Node struct {
	key  map[string]*Node
	lock sync.Mutex  // each node has its own mutex
	data interface{} // if data is not nil, key map should be empty

	pipline *streamcell // allow user to stream data to other location when Node is modified
}

var shardData = map[string]interface{}{"tables": map[string]*table{}}

var test_shardData = map[string]*Node{}

var shardData_lock = sync.Mutex{}

func register_data(loc map[string]interface{}, key string, data interface{}) { // send data to channel

	shardData_lock.Lock() // more testing needed, but adding a lock makes the assignment faster
	loc[key] = data
	shardData_lock.Unlock()
}
