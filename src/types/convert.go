package types

import (
	"dejavuDB/src/types/decimal"
	"dejavuDB/src/types/float128"
	"errors"
	"unsafe"
)

type Value struct {
	Dtype     DType
	Pointer   unsafe.Pointer
	Timestamp uint64
}

func ValueFromBytes(s []byte) (Value, error) {
	v := &Value{}
	err := v.FromBytes(s)
	return *v, err
}

func (val Value) ToBytes() []byte {
	var B []byte
	v := val.Pointer
	switch val.Dtype {
	case String:
		B = []byte(*(*string)(v))
	case Int, Int64:
		b := *(*[8]byte)(v)
		B = b[:]
	case Int32:
	case Int16:
	case Int8:
	case Int128:
	case Uint, Uint64:
		b := *(*[8]byte)(v)
		B = b[:]
	case Uint32:
	case Uint16:
	case Uint8:
	case Uint128:
	case Decimal, Decimal64:
	case Decimal32:
		d := (*decimal.Decimal32)(v)
		a := unsafe.Pointer(&d.I)
		b := unsafe.Pointer(&d.P)
		c := []byte{}
		copy(c, (*(*[2]byte)(a))[:])
		c = append(c, (*(*[2]byte)(b))[:]...)
	case Decimal128:

	case Float, Float64:
		b := *(*[8]byte)(v)
		B = b[:]
	case Float32:
	case Float128:
		B = (*float128.Float128)(v)[:]

	case Byte:
		B = []byte{*(*byte)(v)}
	case Byte_arr:
		B = *(*[]byte)(v)
	case Bool:
		if *(*bool)(v) {
			B = []byte{0x01}
		} else {
			B = []byte{0x00}
		}
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

	return append([]byte{byte(val.Dtype)}, B...)
}

func (v *Value) FromBytes(bs []byte) (err error) {

	if len(bs) == 0 {
		err = errors.New("no bytes provided")
		return
	}
	v.Dtype = DType(bs[0])
	bs = bs[1:]
	switch DType(v.Dtype) {
	case String:
		a := string(bs)
		v.Pointer = unsafe.Pointer(&a)

	case Bool:
		var a bool = false
		if bs[0] > 0x00 {
			a = true
		}
		v.Pointer = unsafe.Pointer(&a)
	case Byte:
		a := bs[0]
		v.Pointer = unsafe.Pointer(&a)
	case Byte_arr:
		v.Pointer = unsafe.Pointer(&bs)

	}
	return
}
