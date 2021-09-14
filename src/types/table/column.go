package table

import (
	"unsafe"

	"../../types"
)

type Column struct {
	Dtype byte
	Data  unsafe.Pointer // pointer to a slice
}

func NewColumn(dtype byte) Column {
	return Column{Dtype: dtype}
}

func (c Column) Add(data interface{}) error {
	switch c.Dtype {
	case types.String:
		*(*[]string)(c.Data) = append(*(*[]string)(c.Data), data.(string))

	case types.Int:
		*(*[]int)(c.Data) = append(*(*[]int)(c.Data), data.(int))

	}
	return nil
}

func (c Column) GetRange(from, end int) ([]unsafe.Pointer, error) {
	var buf []unsafe.Pointer
	return buf, nil
}
