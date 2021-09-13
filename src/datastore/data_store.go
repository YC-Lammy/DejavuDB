package datastore

import (
	"strings"
	"sync"
	"unsafe"
)

var Data = map[string]Node{}

var Data_lock = sync.Mutex{}

var Layers = []Layer{}

type Layer struct {
	Nodes []Node
}

type Node struct {
	name   []byte
	subkey map[string]Node
	lock   *sync.Mutex    // each node has its own mutex
	data   unsafe.Pointer // if data is not nil, key map should be empty
	dtype  byte           // declared at constant
}

func (loc Node) register_data(key string, data interface{}) { // send data to channel
	switch v := data.(type) {
	case Node:
		loc.lock.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key] = v
		loc.lock.Unlock()
	case *Node:
		loc.lock.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key] = *v
		loc.lock.Unlock()

	case unsafe.Pointer:
		loc.data = v

	default:
		loc.data = unsafe.Pointer(&v)
	}
}

func Get(key string) unsafe.Pointer {
	if key == "" {
		return nil
	}
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
	if v, ok := pointer[keys[len(keys)-1]]; ok {
		return v.data
	}
	return nil
}
