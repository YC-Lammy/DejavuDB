package datastore

import (
	"errors"
	"fmt"
	"os"
	"path"
	"src/config"
	"src/types/binjson"
	"src/types/contract"
	"src/types/decimal"
	"src/types/float128"
	"src/types/int128"
	"src/types/uint128"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"src/types"
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

//
//
//
//

func Get(key string) (byte, unsafe.Pointer) {
	if key == "" {
		return 0x00, nil
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
}

//
//
//
//

func Set(key string, data string, dtype byte) error {
	if key == "" {
		return errors.New("invalid empty key")
	}

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

/*
//
//
*/

func (l *Node) write_string_to_loc(data string, dtype byte) (unsafe.Pointer, error) {
	var p unsafe.Pointer

	switch dtype {

	case types.String:
		p = unsafe.Pointer(&data)
		l.lock.RLock()
		l.data = p
		l.lock.RUnlock()

	case types.Int64, types.Int:
		a, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.lock.RLock()
		l.data = p
		l.lock.RUnlock()

	case types.Int32:
		a, err := strconv.ParseInt(data, 10, 32)
		if err != nil {
			return nil, err
		}
		b := int32(a)
		p = unsafe.Pointer(&b)
		l.lock.RLock()
		l.data = p
		l.lock.RUnlock()

	case types.Int16:
		a, err := strconv.ParseInt(data, 10, 16)
		if err != nil {
			return nil, err
		}
		b := int16(a)
		p = unsafe.Pointer(&b)
		l.data = p

	case types.Int8:
		a, err := strconv.ParseInt(data, 10, 8)
		if err != nil {
			return nil, err
		}
		b := int8(a)
		p = unsafe.Pointer(&b)
		l.data = p

	case types.Int128:
		a, err := int128.StrToInt128(data)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Uint, types.Uint64:
		a, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Uint32:
		a, err := strconv.ParseUint(data, 10, 32)
		if err != nil {
			return nil, err
		}
		b := uint32(a)
		p = unsafe.Pointer(&b)
		l.data = p

	case types.Uint16:
		a, err := strconv.ParseUint(data, 10, 16)
		if err != nil {
			return nil, err
		}
		b := uint16(a)
		p = unsafe.Pointer(&b)
		l.data = p

	case types.Uint8:
		a, err := strconv.ParseUint(data, 10, 8)
		if err != nil {
			return nil, err
		}
		b := uint8(a)
		p = unsafe.Pointer(&b)
		l.data = p

	case types.Uint128:
		a, err := uint128.StrToUint128(data)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Decimal, types.Decimal64:
		a, err := decimal.StrToDecimal64(data)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Decimal32:
		a, err := decimal.StrToDecimal32(data)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Decimal128:
		a, err := decimal.StrToDecimal128(data)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Float, types.Float64:
		a, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Float32:
		a, err := strconv.ParseFloat(data, 32)
		if err != nil {
			return nil, err
		}
		b := float32(a)
		p = unsafe.Pointer(&b)
		l.data = p

	case types.Float128:
		a, err := float128.StrToFloat128(data)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Byte:
		a := data[0]
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Byte_arr: // a string
		a := []byte(data)
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Bool:
		a, err := strconv.ParseBool(data)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(&a)
		l.data = p

	case types.Graph:
	case types.Table: // when set table, table is initialized

	case types.Json:
		a, err := binjson.NewBinaryJson([]byte(data))
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(a)
		l.data = p

	case types.SmartContract:
	case types.Contract: // in a json format
		a, err := contract.NewContract(data)
		if err != nil {
			return nil, err
		}
		p = unsafe.Pointer(a)
		l.data = p

	case types.Money:
	case types.SmallMoney:
	case types.Time:
	case types.Date:
	case types.Datetime:
	case types.Smalldatetime:

	case types.Null:
		l.data = nil
		p = nil

	}
	return p, nil
}
