package table

import (
	"errors"
	"os"
	"unsafe"

	"../../types"
)

type Column struct {
	Name  string
	Dtype byte
	Data  unsafe.Pointer // pointer to a slice
}

func NewColumn(dtype byte) Column {
	return Column{Dtype: dtype}
}

func (c *Column) Add(data interface{}) error {
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

func (c *Column) GetRange(from, end int) ([]unsafe.Pointer, error) {
	var buf []unsafe.Pointer
	return buf, nil
}

func (c *Column) ToDisk(path string) error {
	f, err := os.Create(path + string(os.PathSeparator) + c.Name)
	if err != nil {
		return err
	}
	switch c.Dtype {
	case types.String:
		for _, v := range *(*[]string)(c.Data) {
			a, err := types.ToBytes(v)
			if err != nil {
				return err
			}
			f.Write(append(a, '\n'))
		}
	case types.Bool:
	case types.Byte:
	case types.Byte_arr:
	case types.Int, types.Int64:
	}
	return nil
}
