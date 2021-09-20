package table

import (
	"errors"
	"sync"
)

type Table struct {
	Columns    map[string]*Column
	permission [3]uint8
	Leng       uint64
	Lock       sync.RWMutex

	Dtypes []byte

	is_in_ram bool
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

func GetRange(from, to int) map[string]Column {

}

func (t *Table) Length() uint64 {
	return t.Leng
}

func (t *Table) Wideth() uint64 {
	return uint64(len(t.Columns))
}
