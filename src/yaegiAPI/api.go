package yaegiAPI

import (
	"src/datastore"
)

type db interface {
	Set(key string, data interface{}, dtype byte) error
	Get(key string) value
	Update(key string, data interface{})
	Move(string, string) error
	Types() types_struct
}

type types_struct struct {
	String,
	Int, Int64, Int32, Int16, Int8, Int128,
	Uint, Uint64, Uint32, Uint16, Uint8, Uint128,
	Float, Float64, Float32, Float128 byte
}

type database struct {
	uid uint64
	gid uint64
}

func (d *database) Set(key string, data interface{}, dtype byte) error {
	return nil
}

func (d *database) Get(key string) value {
	dtype, ptr := datastore.Get(key)
	return value{Ptr: ptr, Dtype: dtype}
}

func (d *database) Update(key string, data interface{}) {}

func (d *database) Move(src, dst string) error {
	return nil
}

func (d *database) Types() types_struct {
	return types_struct{
		String: 0x00,
	}
}
