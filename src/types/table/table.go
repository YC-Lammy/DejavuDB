package table

import "errors"

type Table struct {
	Columns    map[string]Column
	permission [3]uint8
	length     uint64
}

func (t Table) Insert(names []string, data []interface{}) error {
	a := len(data)
	if a != len(t.Columns) {
		return errors.New("values and columns not match")
	}
	for i := 0; i < a; i++ {
		t.Columns[names[i]].Add(data[i])
	}
	return nil
}

func (t Table) Length() uint64 {
	return t.length
}

func (t Table) Wideth() uint64 {
	return uint64(len(t.Columns))
}
