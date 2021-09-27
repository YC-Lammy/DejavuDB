package table

import (
	"errors"
	"os"
	"unsafe"

	"src/types"
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
	case types.Int32:
	case types.Int16:
	case types.Int8:
	case types.Int128:
	case types.Uint, types.Uint64:
	case types.Uint32:
	case types.Uint16:
	case types.Uint8:
	case types.Uint128:
	case types.Decimal, types.Decimal64:
	case types.Decimal32:
	case types.Decimal128:
	case types.Float, types.Float64:
	case types.Float32:
	case types.Float128:
	case types.Byte:
	case types.Byte_arr:
	case types.Bool:
	case types.Graph:
	case types.Table:
	case types.Json:
	case types.SmartContract:
	case types.Contract:
	case types.Money:
	case types.SmallMoney:
	case types.Time:
	case types.Date:
	case types.Datetime:
	case types.Smalldatetime:
	case types.Null:

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
	bytes := []byte{}
	switch c.Dtype {
	case types.String:
		for _, v := range *(*[]string)(c.Data) {
			a, err := types.ToBytes(v)
			if err != nil {
				return err
			}
			bytes = append(bytes, append(a, '\n')...)
		}
	case types.Int, types.Int64:
	case types.Int32:
	case types.Int16:
	case types.Int8:
	case types.Int128:
	case types.Uint, types.Uint64:
	case types.Uint32:
	case types.Uint16:
	case types.Uint8:
	case types.Uint128:
	case types.Decimal, types.Decimal64:
	case types.Decimal32:
	case types.Decimal128:
	case types.Float, types.Float64:
	case types.Float32:
	case types.Float128:
	case types.Byte:
	case types.Byte_arr:
	case types.Bool:
	case types.Graph:
	case types.Table:
	case types.Json:
	case types.SmartContract:
	case types.Contract:
	case types.Money:
	case types.SmallMoney:
	case types.Time:
	case types.Date:
	case types.Datetime:
	case types.Smalldatetime:
	case types.Null:
	}
	f.Write(bytes)
	return nil
}
