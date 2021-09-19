package table

import (
	"errors"
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
		if v, ok := data.(string); ok {
			*(*[]string)(c.Data) = append(*(*[]string)(c.Data), v)
		} else {
			return errors.New("type not match")
		}
	case types.Int, types.Int64:
		if v, ok := data.(int64); ok {
			*(*[]int64)(c.Data) = append(*(*[]int64)(c.Data), v)
		} else {
			return errors.New("type not match")
		}
	}
	return nil
}

func (c Column) GetRange(from, end int) ([]unsafe.Pointer, error) {
	var buf []unsafe.Pointer
	return buf, nil
}
