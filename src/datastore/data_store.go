package datastore

import (
	"strings"
	"sync"
	"time"
)

const (
	Byte = 0x00

	Int = 0x01
	Int8 = 0x02
	Int16 = 0x03
	Int32 = 0x04
	Int64 = 0x05
	Int128 = 0x06

	Float = 0x07
	Float32 = 0x08
	Float64 = 0x09
	Float128 = 0x10

	Decimal32 = 0x11
	Decimal64 = 0x12
	Decimal128 = 0x13

	String = 0x14
	Bool = 0x15
	Table = 0x16

)
var Data = map[string]Node{}

var Data_lock = sync.Mutex{}

var Layers = []Layer{}

type Layer struct{
	Nodes []Node
}

type Node struct {
	name        string
	subkey      map[string]*Node
	lock        sync.Mutex  // each node has its own mutex
	data        interface{} // if data is not nil, key map should be empty
	dtype       byte // declared at constant
}

func (loc Node) register_data(key string, data interface{}) { // send data to channel
	switch v:= data.(type){
	case Node:
		loc.lock.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key] = v
		loc.lock.Unlock()
	case *Node:
		loc.lock.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key] = *v
		loc.lock.Unlock()

	default:
		loc.data = data
	}
}

func Get(key string, json_options ...map[string]interface{}) interface{} {
	if key == ""{
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
	if v, ok := pointer[keys[len(keys)-1]];ok{
		return v.data
	}
	return nil
}
