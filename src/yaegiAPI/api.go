package yaegiAPI

import (
	"unsafe"
)

type db interface {
	Set(key string, data interface{}, dtype byte) error
	Get(key string)
	Update(key string, data interface{})
	Move(string, string) error
	Types() types

	WriteError(error)
}

type types struct {
	String,
	Int, Int64, Int32, Int16, Int8, Int128,
	Uint, Uint64, Uint32, Uint16, Uint8, Uint128,
	Float, Float64, Float32, Float128 byte
}

type value struct {
	Ptr   unsafe.Pointer
	Dtype byte
}

type database struct {
	uid uint64
	gid uint64
}

func (d *database) Set(key string, data interface{}, dtype byte) error {
	return nil
}

func (d *database) Get(key string)

func (d *database) Update(key string, data interface{}) {}

func (d *database) Move(src, dst string) error {
	return nil
}

func (d *database) Types() types {
	return types{
		String: 0x00,
	}
}

func (d *database) WriteError(err error) {

}
