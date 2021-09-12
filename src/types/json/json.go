package typejson

import (
	"rogchip.com/v8go"

	"../../lazy"
)

var vm *v8go.Context

type json struct {
	location string
}

func NewJson(str string) (*json, error) {
	loc := lazy.RandString(8)
	vm.RunScript(loc+"="+str, "json.js")
	return &json{location: loc}, nil
}

func (j json) Delete() {
	vm.RunScript("delete "+j.location, "delete.js")
}

func (j json) String() string {
	val, _ := vm.RunScript(j.location, "get.js")
	return val.String()
}
