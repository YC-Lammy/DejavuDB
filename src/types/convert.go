package types

import (
	"errors"
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
