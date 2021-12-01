package datastore

import (
	"dejavuDB/src/types"
	"errors"
	"os"
	"path"
	"strings"
	"sync"
	"unsafe"
)

var (
	InvalidKeyError = errors.New("invalid key")
)

var userhomedir, _ = os.UserHomeDir()
var database_path = path.Join(userhomedir, "dejavuDB", "database")

var Data = map[string]*Node{}

var Data_lock = sync.RWMutex{}

type Node struct {
	subkey map[string]*Node
	lock   sync.RWMutex // each node has its own mutex
	//data_lock sync.Mutex
	data  unsafe.Pointer
	dtype byte // declared at types
}

/*
func (loc *Node) register_data(data interface{}, dtype byte, key string) error { // send data to channel
	var t unsafe.Pointer
	switch v := data.(type) {
	case Node:

		loc.lock.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key] = &v
		loc.lock.Unlock()
		return nil
	case *Node:

		loc.lock.Lock() // more testing needed, but adding a lock makes the assignment faster
		loc.subkey[key] = v
		loc.lock.Unlock()
		return nil

	case unsafe.Pointer:
		loc.lock.RLock()
		loc.data = v
		loc.lock.RUnlock()

		f, err := os.Create(path.Join(database_path, key))
		if err != nil {
			return err
		}
		b, err := types.ToBytes(v, dtype)
		if err != nil {
			return err
		}
		f.Write(b)
		f.Close()

	case string: // a javascript string from value
		ptr, err := loc.write_string_to_loc(v, dtype)
		if err != nil {
			return err
		}
		t = ptr

	case int:
		if dtype != types.Int && dtype != types.Int64 {
			return errors.New("register: unexpected type int")
		}
		a := int64(v)
		b := unsafe.Pointer(&a)
		loc.lock.Lock()
		loc.data = b
		loc.dtype = dtype
		loc.lock.Unlock()
		t = b

	default:
		return errors.New("register: unsupported type")

	}
	f, err := os.Create(path.Join(database_path, key))
	if err != nil {
		return err
	}
	b, err := types.ToBytes(t, dtype)
	if err != nil {
		return err
	}
	f.Write(b)
	f.Close()
	return nil
}
*/

//
//
//
//

func Get(key string) types.Value {
	/*
		if key == "" {
			return nil
		}
		if config.Debug {
			fmt.Println("Get key " + key)
		}

		Data_lock.RLock()
		var pointer = Data // copy pointers into steak
		Data_lock.RUnlock()

		keys := strings.Split(key, ".")
		if len(keys) == 1 {
			if v, ok := pointer[keys[0]]; ok {
				v.lock.RLock()
				a, b := v.dtype, v.data
				v.lock.RUnlock()
				return a, b
			}
			return 0x00, nil
		}
		for _, v := range keys[0 : len(keys)-1] {
			if v, ok := pointer[v]; ok {
				v.lock.RLock()
				pointer = v.subkey
				v.lock.RUnlock()
			} else {
				return 0x00, nil
			}
		}
		if v, ok := pointer[keys[len(keys)-1]]; ok {
			v.lock.RLock()
			a, b := v.dtype, v.data
			v.lock.RUnlock()
			return a, b
		}
		return 0x00, nil
	*/
	f := strings.Split(key, ".")
	switch len(f) {
	case 1:
		if v, ok := layer1[key]; ok {
			return v
		}
	case 2:
		if v, ok := layer2[key]; ok {
			return v
		}
	case 3:
		if v, ok := layer3[key]; ok {
			return v
		}
	case 4:
		if v, ok := layer4[key]; ok {
			return v
		}
	case 5:
		if v, ok := layer5[key]; ok {
			return v
		}
	case 6:
		if v, ok := layer6[key]; ok {
			return v
		}
	case 7:
		if v, ok := layer7[key]; ok {
			return v
		}
	case 8:
		if v, ok := layer8[key]; ok {
			return v
		}
	case 9:
		if v, ok := layer9[key]; ok {
			return v
		}
	case 10:
		if v, ok := layer10[key]; ok {
			return v
		}
	}
	return types.Value{
		Dtype: types.Null,
	}
}

//
//
//
//

func Set(key string, val types.Value) error {
	if key == "" {
		return errors.New("invalid empty key")
	}
	/*

		var pointer = Data // copy pointers into steak

		keys := strings.Split(key, ".")

		if len(keys) == 1 { // only one key provide
			if n, ok := pointer[keys[0]]; ok {
				n.register_data(data, dtype, key)
				return nil
			} else {
				n := CreateNode()
				pointer[keys[0]] = n
				n.register_data(data, dtype, key)
				return nil
			}
		}

		for _, v := range keys[0 : len(keys)-1] {
			if n, ok := pointer[v]; ok {
				pointer = n.subkey
			} else {
				n := CreateNode()
				pointer[v] = n
				n.register_data(data, dtype, key)
				return nil
			}
		}
		if n, ok := pointer[keys[len(keys)-1]]; ok {
			n.register_data(data, dtype, key)

		} else {
			n := CreateNode()
			pointer[keys[len(keys)-1]] = n
			n.register_data(data, dtype, key)
		}
	*/

	f := strings.Split(key, ".")

	switch len(f) {
	case 1:
		layer1[key] = val
	case 2:
		layer2[key] = val
	case 3:
		layer3[key] = val
	case 4:
		layer4[key] = val
	case 5:
		layer5[key] = val
	case 6:
		layer6[key] = val
	case 7:
		layer7[key] = val
	case 8:
		layer8[key] = val
	case 9:
		layer9[key] = val
	case 10:
		layer10[key] = val
	}
	return nil
}

func Delete(key string) error {
	if key == "" {
		return errors.New("invalid empty key")
	}

	f := strings.Split(key, ".")
	switch len(f) {
	case 1:
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	case 7:
	case 8:
	case 9:
	case 10:
	}
	return nil
}

//
//
//
//

func CreateNode() *Node {
	return &Node{subkey: map[string]*Node{}, lock: sync.RWMutex{}}
}

//
//
//
//

/*
func Update(key string, data interface{}, dtype ...byte) error {
	var node *Node
	if key == "" {
		return errors.New("invalid key")
	}

	Data_lock.RLock()
	var pointer = Data // copy pointers into steak
	Data_lock.RUnlock()

	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		if v, ok := pointer[keys[0]]; ok {
			node = v
		}
		return InvalidKeyError
	}
	for _, v := range keys[0 : len(keys)-1] {
		if v, ok := pointer[v]; ok {
			v.lock.RLock()
			pointer = v.subkey
			v.lock.RUnlock()
		} else {
			return InvalidKeyError
		}
	}
	if v, ok := pointer[keys[len(keys)-1]]; ok {
		node = v
	} else {
		return InvalidKeyError
	}

	node.lock.RLock()
	t := node.dtype
	node.lock.RUnlock()

	switch data.(type) {
	case unsafe.Pointer:
		if dtype[0] != t {
			return errors.New("Update: type mismatch")
		}
		return node.register_data(data, dtype[0], key)
	}

	return node.register_data(data, t, key)
}
*/

/*
//
//
*/
