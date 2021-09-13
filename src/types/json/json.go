package typejson

import (
	"rogchip.com/v8go"

	"../../lazy"
)

var vm *v8go.Context

type json struct {
	location [8]byte
}

func NewJson(str string) (*json, error) {
	s := lazy.RandString(8)
	vm.RunScript(s+"="+str, "json.js")
	var loc [8]byte
	copy(loc[:], s)
	return &json{location: loc}, nil
}

func (j json) Delete() {
	vm.RunScript("delete "+string(j.location[:]), "delete.js")
}

func (j json) String() string {
	val, _ := vm.RunScript(j.location, "get.js")
	return val.String()
}
