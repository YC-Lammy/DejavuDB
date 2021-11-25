package types

import (
	"dejavuDB/src/types/float128"
	"dejavuDB/src/types/int128"
	"dejavuDB/src/types/uint128"
	"unsafe"
)

type embedInterface struct {
	typ unsafe.Pointer
	_   unsafe.Pointer
}

var TypePtrRegister = map[unsafe.Pointer]DType{}

func init() {
	l := map[interface{}]DType{
		"hello":             String,
		int(0):              Int,
		int64(0):            Int64,
		int32(0):            Int32,
		int16(0):            Int16,
		int8(0):             Int8,
		int128.Int128{}:     Int128,
		uint(0):             Uint,
		uint64(0):           Uint64,
		uint32(0):           Uint32,
		uint16(0):           Uint16,
		uint8(0):            Uint8,
		uint128.Uint128{}:   Uint128,
		float64(0):          Float64,
		float32(0):          Float32,
		float128.Float128{}: Float128,
	}
	for k, v := range l {
		TypePtrRegister[(*embedInterface)(unsafe.Pointer(&k)).typ] = v
	}
}
