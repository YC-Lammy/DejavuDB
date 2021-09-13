package table

import "unsafe"

type Column struct {
	Dtype byte
	Data  []unsafe.Pointer
}

func (column Column) Add(data interface{}) error {
	switch column.Dtype {

	}
	return nil
}

func (c Column) GetRange(from, end int) ([]unsafe.Pointer, error) {
	var buf []unsafe.Pointer
	return buf, nil
}
