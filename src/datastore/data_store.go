package datastore

import (
	"errors"
	"strings"
	"sync"
	"unsafe"

	"../types"
)

var Data = map[string]*Node{}

var Data_lock = sync.Mutex{}

type Node struct {
	subkey map[string]*Node
	lock   sync.Mutex // each node has its own mutex
	//data_lock sync.Mutex
	data  unsafe.Pointer
	dtype byte // declared at types
}

func (loc *Node) register_data(data interface{}, key ...string) { // send data to channel
	switch v := data.(type) {
	case Node:
		l := &loc.lock
		l.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key[0]] = &v
		l.Unlock()
	case *Node:
		l := &loc.lock
		l.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key[0]] = v
		l.Unlock()

	case unsafe.Pointer:
		//loc.data_lock.Lock()
		loc.data = v
		//loc.data_lock.Unlock()

	case string: // a javascript string from value
		write_type_to_loc(loc, v, key[0])

	default:

	}
}

func Get(key string) unsafe.Pointer {
	if key == "" {
		return nil
	}

	var pointer = Data // copy pointers into steak

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
		} else {
			return nil
		}
	}
	if v, ok := pointer[keys[len(keys)-1]]; ok {
		return v.data
	}
	return nil
}

func Set(key string, data string, dtype byte) error {
	if key == "" {
		return errors.New("invalid empty key")
	}

	var pointer = Data // copy pointers into steak

	keys := strings.Split(key, ".")

	if len(keys) == 1 { // only one key provide
		if v, ok := pointer[keys[0]]; ok {
			v.register_data(data, string(dtype))
		}
		return nil
	}

	for _, v := range keys[0 : len(keys)-1] {
		if v, ok := pointer[v]; ok {
			pointer = v.subkey
		} else {
			break
		}
	}
	if v, ok := pointer[keys[len(keys)-1]]; ok {
		v.register_data(data, string(dtype))

		return nil
	}
	return nil
}

func write_type_to_loc(l *Node, data string, dtype string) error {
	switch dtype[0] {

	case types.Table:

	case types.String:

	case types.Byte:
	case types.Contract:
	case types.SmartContract:

	}
	return nil
}
