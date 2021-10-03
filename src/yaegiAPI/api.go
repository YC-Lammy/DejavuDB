package yaegiAPI

import (
	"errors"
	"src/datastore"
)

type db interface {
	Set(key string, data interface{}, dtype byte) error
	Get(key string) value
	Update(key string, data interface{}) error
	Move(string, string) error
	Types() types_struct
}

type types_struct struct {
	String,
	Int, Int64, Int32, Int16, Int8, Int128,
	Uint, Uint64, Uint32, Uint16, Uint8, Uint128,
	Float, Float64, Float32, Float128 byte
}

type Database struct {
	uid       uint64
	gid       uint64
	reversefn []func()
}

func (d *Database) Set(key string, data interface{}, dtype byte) error {
	switch v := data.(type) {
	case value:
		if v.(Val).Dtype != dtype {
			return errors.New("Set: type mismatch")
		}
		fn, err := datastore.JsSet(key, v.(Val).Ptr, v.(Val).Dtype)
		if err != nil {
			return err
		}
		d.reversefn = append(d.reversefn, fn)
	case string:
	}
	return nil
}

func (d *Database) Get(key string) Val {
	dtype, ptr := datastore.Get(key)
	return Val{Ptr: ptr, Dtype: dtype}
}

func (d *Database) Update(key string, data interface{}) error {
	switch v := data.(type) {
	case value:
		s := v.(Val)
		err := datastore.Update(key, s.Ptr, s.Dtype)
		if err != nil {
			return err
		}
		//d.reversefn = append(d.reversefn, fn)
	}
	return nil
}

func (d *Database) Move(src, dst string) error {
	return nil
}

func (d *Database) Types() types_struct {
	return types_struct{
		String: 0x00,
	}
}
