package types

import (
	"errors"
	"unsafe"
)

func ToBytes(data unsafe.Pointer, dtype byte) ([]byte, error) {
	var B []byte
	v := data
	switch dtype {
	case String:
		return ToBytes(*(*string)(v), dtype...)
	case Int, Int64:
		b := *(*[8]byte)(v)
		B = b[:]
	case Int32:
	case Int16:
	case Int8:
	case Int128:
	case Uint, Uint64:
	case Uint32:
	case Uint16:
	case Uint8:
	case Uint128:
	case Decimal, Decimal64:
	case Decimal32:
		a := unsafe.Pointer(&v.I)
		b := unsafe.Pointer(&v.P)
		c := []byte{}
		copy(c, (*(*[2]byte)(a))[:])
		c = append(c, (*(*[2]byte)(b))...)
	case Decimal128:
	case Float, Float64:
	case Float32:
	case Float128:
		B = v[:]
	case Byte:
	case Byte_arr:
		return *(*[]byte)(v), nil
	case Bool:
		return ToBytes(*(*bool)(v), dtype...)
	case Graph:
	case Table:
	case Json:
	case SmartContract:
	case Contract:
	case Money:
	case SmallMoney:
	case Time:
	case Date:
	case Datetime:
	case Smalldatetime:
	case Null:
	}

	return append([]byte{dtype}, B...), nil
}

func FromBytes(bs []byte) (p unsafe.Pointer, dtype byte, err error) {

	if len(bs) == 0 {
		err = errors.New("no bytes provided")
		return
	}
	dtype = bs[0]
	bs = bs[1:]
	switch dtype {
	case String:
		a := string(bs)
		p = unsafe.Pointer(&a)

	case Bool:
		var a bool = false
		if bs[0] > 0x00 {
			a = true
		}
		p = unsafe.Pointer(&a)
	case Byte:
		a := bs[0]
		p = unsafe.Pointer(&a)
	case Byte_arr:
		p = unsafe.Pointer(&bs)

	}
	return
}
