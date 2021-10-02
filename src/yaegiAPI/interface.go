package yaegiAPI

var init_script = `
import (
	"fmt"
	"log"
	"math"
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

type value interface{
	String() string
	Interface() interface{}
	Add() error
	Sub() error
}

func main(DB db) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

`
var ending = `
	return nil
}
`
