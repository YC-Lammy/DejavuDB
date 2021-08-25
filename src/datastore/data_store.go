package datastore

import (
	"strings"
	"sync"
	"time"
)

var Data = map[string]*Node{}

var Data_lock = sync.Mutex{}

type Node struct {
	name        string
	subkey      map[string]*Node
	lock        sync.Mutex  // each node has its own mutex
	data        interface{} // if data is not nil, key map should be empty
	create_time time.Time
	modify_time time.Time

	//pipline_cell   *streamcell // allow user to stream data to other location when Node is modified
	trigger_script string // trigger java script when event happened
}

func (loc *Node) register_data(key string, data interface{}) { // send data to channel

	if v, ok := data.(*Node); ok {
		loc.lock.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key] = v
		loc.lock.Unlock()
	} else {
		loc.data = data
	}
}

func Get(key string, json_options ...map[string]interface{}) interface{} {
	Data_lock.Lock()
	var pointer = Data
	Data_lock.Unlock()
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		if v, ok := pointer[keys[0]]; ok {
			return v.data
		}
		return nil
	}
	for _, v := range keys[0 : len(keys)-1] {
		if v, ok := pointer[v]; ok {
			pointer = v.subkey
		}
	}
	return nil
}
