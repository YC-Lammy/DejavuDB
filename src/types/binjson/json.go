package binjson

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"../../types"
)

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
		r += GetElemStr(b[i+1])
		i += 2
	}
	r += "}"
	return r
}

func GetElemStr(b []byte) string {
	var a string
	switch b[0] {
	case types.Bool:
		switch b[1] {
		case 0x00:
			a = "false"
		case 0x01:
			a = "true"
		}
	case types.Float64:
		var length float64
		buf := bytes.NewReader(b[1:])
		binary.Read(buf, binary.LittleEndian, &length)
		a = strconv.FormatFloat(length, 'f', 15, 64)
	case types.String:
		a = "'" + string(b[1:]) + "'"
	case types.Array_interface: // array
		// each elem separated by ";%@"

		buf := [][]byte{}
		sepc := 0
		bufc := 0

		for _, i := range b[1:] {
			switch i {
			case ';':
				if sepc != 0 {
					sepc = 0
					buf[bufc] = append(buf[bufc], i)
				} else {
					sepc += 1
				}

			case '%':
				if sepc != 1 {
					sepc = 0
					buf[bufc] = append(buf[bufc], i)
				} else {
					sepc += 1
				}
			case '@':
				if sepc != 2 {
					sepc = 0
					buf[bufc] = append(buf[bufc], i)
				} else {
					bufc += 1
					sepc = 0
				}

			default:
				buf[bufc] = append(buf[bufc], i)
			}
		}
		a += "["

		for _, v := range buf {
			a += GetElemStr(v) + ","
		}
		a += "]"

	case types.Map_string_interface: // json
		// each elem separated by ";%@"

		buf := [][]byte{}
		sepc := 0
		bufc := 0

		for _, i := range b[1:] {
			switch i {
			case ';':
				if sepc != 0 {
					sepc = 0
					buf[bufc] = append(buf[bufc], i)
				} else {
					sepc += 1
				}

			case '%':
				if sepc != 1 {
					sepc = 0
					buf[bufc] = append(buf[bufc], i)
				} else {
					sepc += 1
				}
			case '@':
				if sepc != 2 {
					sepc = 0
					buf[bufc] = append(buf[bufc], i)
				} else {
					bufc += 1
					sepc = 0
				}

			default:
				buf[bufc] = append(buf[bufc], i)
			}
		}
		a += "{"
		l := len(buf)
		i := 0
		for i < l {
			a += "'" + string(buf[i]) + "':"
			a += GetElemStr(buf[i+1])
			i += 2
		}
		a += "}"

	case types.Null:
		a = "null"
	}
	return a
}
