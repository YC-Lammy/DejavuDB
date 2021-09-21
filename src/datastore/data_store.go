package datastore

import (
	"errors"
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
		loc.write_type_to_loc(v, key[0])

	default:

	}
}

func Get(key string) (byte, unsafe.Pointer) {
	if key == "" {
		return 0x00, nil
	}

	var pointer = Data // copy pointers into steak

	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		if v, ok := pointer[keys[0]]; ok {
			return v.dtype, v.data
		}
		return 0x00, nil
	}
	for _, v := range keys[0 : len(keys)-1] {
		if v, ok := pointer[v]; ok {
			pointer = v.subkey
		} else {
			return 0x00, nil
		}
	}
	if v, ok := pointer[keys[len(keys)-1]]; ok {
		return v.dtype, v.data
	}
	return 0x00, nil
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

func (l *Node) write_type_to_loc(data string, dtype string) error {

	switch dtype[0] {

	case types.String:
		l.data = unsafe.Pointer(&data)

	case types.Int64, types.Int:
		a, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Int32:
		a, err := strconv.ParseInt(data, 10, 32)
		if err != nil {
			return err
		}
		b := int32(a)
		l.data = unsafe.Pointer(&b)

	case types.Int16:
		a, err := strconv.ParseInt(data, 10, 16)
		if err != nil {
			return err
		}
		b := int16(a)
		l.data = unsafe.Pointer(&b)

	case types.Int8:
		a, err := strconv.ParseInt(data, 10, 8)
		if err != nil {
			return err
		}
		b := int8(a)
		l.data = unsafe.Pointer(&b)

	case types.Int128:
		a, err := int128.StrToInt128(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Uint, types.Uint64:
		a, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Uint32:
		a, err := strconv.ParseUint(data, 10, 32)
		if err != nil {
			return err
		}
		b := uint32(a)
		l.data = unsafe.Pointer(&b)

	case types.Uint16:
		a, err := strconv.ParseUint(data, 10, 16)
		if err != nil {
			return err
		}
		b := uint16(a)
		l.data = unsafe.Pointer(&b)

	case types.Uint8:
		a, err := strconv.ParseUint(data, 10, 8)
		if err != nil {
			return err
		}
		b := uint8(a)
		l.data = unsafe.Pointer(&b)

	case types.Uint128:
		a, err := uint128.StrToUint128(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Decimal, types.Decimal64:
		a, err := decimal.StrToDecimal64(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Decimal32:
		a, err := decimal.StrToDecimal32(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Decimal128:
		a, err := decimal.StrToDecimal128(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Float, types.Float64:
		a, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Float32:
		a, err := strconv.ParseFloat(data, 32)
		if err != nil {
			return err
		}
		b := float32(a)
		l.data = unsafe.Pointer(&b)

	case types.Float128:
		a, err := float128.StrToFloat128(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)

	case types.Byte:
		a := data[0]
		l.data = unsafe.Pointer(&a)

	case types.Byte_arr: // a string
		a := []byte(data)
		l.data = unsafe.Pointer(&a)

	case types.Bool:
		a, err := strconv.ParseBool(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(&a)
	case types.Graph:
	case types.Table: // when set table, table is initialized

	case types.Json:
		a, err := binjson.NewBinaryJson(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(a)

	case types.SmartContract:
	case types.Contract: // in a json format
		a, err := contract.NewContract(data)
		if err != nil {
			return err
		}
		l.data = unsafe.Pointer(a)

	case types.Money:
	case types.SmallMoney:
	case types.Time:
	case types.Date:
	case types.Datetime:
	case types.Smalldatetime:

	case types.Null:
		l.data = nil

	}
	return nil
}
