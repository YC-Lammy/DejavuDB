package typejson

import (
	"rogchip.com/v8go"

	"../../lazy"
	"../../types"
)

var vm *v8go.Context

type BinaryJson [][]byte

func NewBinaryJson(str string) (*BinaryJson, error) {
	return nil, nil
}

func (b BinaryJson) GetElem(key string) []byte {
	a := len(b)
	i := 0
	for i < a {
		if string(b[i]) == key {
			return b[i+1]
		}
		i += 2
	}
	return nil
}

func (b BinaryJson) String() string {
	var r = "{"
	l := len(b)
	i := 0
	for i < l {
		r += "'" + string(b[i]) + "':"
		r += string(b[i+1])
		i += 2
	}
	return ""
}

func (b BinaryJson) GetElemStrByIndex(i int) string {
	var a string
	switch b[i][0] {
	case types.Bool:
		switch b[i][1] {
		case 0x00:
			a = "false"
		case 0x01:
			a = "true"
		}
	case types.Float64:
	case types.String:
	case types.Array_interface: // array
	case types.Map_string_interface: // json
	case types.Null:
		a = "null"
	}
	return a
}

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
