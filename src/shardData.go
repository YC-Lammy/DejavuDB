package main

import (
	"sync"
	"time"
)

// ====================================================
const ( // node events
	node_update = "update"
	node_set    = "set"
	node_get    = "get"
)

//=====================================================

type Node struct {
	name        string
	key         map[string]*Node
	lock        sync.Mutex  // each node has its own mutex
	data        interface{} // if data is not nil, key map should be empty
	create_time time.Time
	modify_time time.Time

	pipline_cell   *streamcell // allow user to stream data to other location when Node is modified
	trigger_script string      // trigger java script when event happened
}

var shardData = map[string]interface{}{"tables": map[string]*table{}} // default table entry

var test_shardData = map[string]*Node{"tables": &Node{key: map[string]*Node{}, lock: sync.Mutex{}}}

var shardData_lock = sync.Mutex{}

func register_data(loc map[string]interface{}, key string, data interface{}) { // send data to channel

	shardData_lock.Lock() // more testing needed, but adding a lock makes the assignment faster
	loc[key] = data
	shardData_lock.Unlock()
}
