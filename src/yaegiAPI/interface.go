package yaegiAPI

var init_script = `
import (
	"fmt"
	"log"
	"math"
)
type db interface {
	Set(key string, data interface{}) error
	Get(key string)
	Update(key, data interface{})
	Move(string, string)
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
type int128 [16]byte

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
