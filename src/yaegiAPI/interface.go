package yaegiAPI

var init_script = `
package foo
import (
	"fmt"
	"log"
	"math"
	"unsafe"
)

type value interface {
	String() string
	Interface() interface{}
	Add(interface{}) error
	Sub(interface{}) error
}
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


func m(DB db) (err error) {
	//DB := *(*db)(huijk)

`
var ending = `
	return nil
}
`
