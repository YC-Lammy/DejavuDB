package types

import (
	"errors"
	"src/types/decimal"
	"src/types/float128"
	"src/types/int128"
	"strings"
	"unsafe"
)

func ToBytes(data interface{}, dtype ...byte) ([]byte, error) {
	switch v := data.(type) {
	case unsafe.Pointer:
		switch dtype[0] {
		case String:
		case Bool:
		}
	case string:
		v = strings.ReplaceAll(v, "\n", `\n`)
		return []byte(v), nil
	case bool:
		if v {
			return []byte{0x01}, nil
		}
		return []byte{0x00}, nil

	case []byte:
		return v, nil

	case int64:
		a := unsafe.Pointer(&v)
		b := *(*[8]byte)(a)
		return b[:], nil
	case int32:
		a := unsafe.Pointer(&v)
		b := *(*[4]byte)(a)
		return b[:], nil
	case int16:
		a := unsafe.Pointer(&v)
		b := *(*[2]byte)(a)
		return b[:], nil
	case int8:
		a := unsafe.Pointer(&v)
		b := *(*[1]byte)(a)
		return b[:], nil
	case int128.Int128:
		return v[:], nil

	case uint64:
		a := unsafe.Pointer(&v)
		b := *(*[8]byte)(a)
		return b[:], nil
	case uint32:
		a := unsafe.Pointer(&v)
		b := *(*[4]byte)(a)
		return b[:], nil
	case uint16:
		a := unsafe.Pointer(&v)
		b := *(*[2]byte)(a)
		return b[:], nil
	case uint8:
		a := unsafe.Pointer(&v)
		b := *(*[1]byte)(a)
		return b[:], nil

	case float64:
		a := unsafe.Pointer(&v)
		b := *(*[8]byte)(a)
		return b[:], nil

	case float32:
		a := unsafe.Pointer(&v)
		b := *(*[4]byte)(a)
		return b[:], nil

	case float128.Float128:
		return v[:], nil

	case decimal.Decimal32:
		a := unsafe.Pointer(&v.I)
		b := unsafe.Pointer(&v.P)
		c := []byte{}
		copy(c, (*(*[2]byte)(a))[:])
		c = append(c, (*(*[2]byte)(b))...)
		return c, nil
	}

	return nil, nil
}

func FromBytes(dtype byte, bs []byte) (p unsafe.Pointer, err error) {

	if len(bs) == 0 && dtype != Null {
		err = errors.New("no bytes provided")
		return
	}
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
