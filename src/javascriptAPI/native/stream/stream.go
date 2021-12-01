package stream

import (
	"io"

	"github.com/dop251/goja"
)

type Duplex struct {
	this *goja.Object
	io   io.ReadWriter
}

func NewDuplex(vm *goja.Runtime, rw io.ReadWriter) *goja.Object {
	d := &Duplex{}
	d.this = vm.NewDynamicObject(d)
	return d.this
}

func (s *Duplex) Get(key string) goja.Value {
	switch key {
	case "write":
	case "read":
	}
	return goja.Undefined()
}
