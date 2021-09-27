package table

import (
	"errors"
	"src/types"
	"src/types/binjson"
	"src/types/decimal"
	"src/types/float128"
	"src/types/int128"
	"src/types/uint128"
	"sync"
	"unsafe"
)

type Table struct {
	Columns    map[string]*Column
	permission [3]uint8
	Leng       uint64
	Lock       sync.RWMutex

	Dtypes []byte

	is_in_ram bool
}

func (t *Table) Length() uint64 {
	return t.Leng
}

func (t *Table) Wideth() uint64 {
	return uint64(len(t.Columns))
}

func (t *Table) Insert(names []string, data []interface{}) error {
	a := len(data)
	if a != len(t.Columns) {
		return errors.New("values and columns not match")
	}
	t.Lock.Lock()
	defer t.Lock.Unlock()
	for i := 0; i < a; i++ {
		t.Columns[names[i]].Add(data[i])
	}
	t.Leng += 1
	return nil
}

func GetRange(from, to int) []Column {

}

func (t *Table) AddColumn(name string, dtype byte) {
	c := Columm{
		Name:  name,
		Dtype: dtype,
	}
	var u unsafe.Pointer
	switch dtype {
	case types.String:
		a := make([]string, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Int, types.Int64:
		a := make([]int64, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Int32:
		a := make([]int32, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Int16:
		a := make([]int16, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Int8:
		a := make([]int8, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Int128:
		a := make([]int128.Int128, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Uint, types.Uint64:
		a := make([]uint64, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Uint32:
		a := make([]uint32, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Uint16:
		a := make([]uint16, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Uint8:
		a := make([]uint8, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Uint128:
		a := make([]uint128.Uint128, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Decimal, types.Decimal64:
		a := make([]decimal.Decimal64, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Decimal32:
		a := make([]decimal.Decimal32, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Decimal128:
		a := make([]decimal.Decimal128, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Float, types.Float64:
		a := make([]float64, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Float32:
		a := make([]float32, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Float128:
		a := make([]float128.Float128, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Byte:
		a := make([]byte, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Byte_arr:
		a := make([][]byte, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Bool:
		a := make([]bool, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Json:
		a := make([]binjson.Json, t.Leng)
		u = unsafe.Pointer(&a)
	case types.Money:
	case types.SmallMoney:
	case types.Time:
	case types.Date:
	case types.Datetime:
	case types.Smalldatetime:
	case types.Null:

	}
	c.Data = u
	t.Columns[name] = &c
}
