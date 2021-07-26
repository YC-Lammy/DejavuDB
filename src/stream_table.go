package main

import (
	"sync"
)

type streamTable struct {
	cells      []*streamcell
	shape      int
	entrycount int

	destination interface{}
}

type streamcell struct { // cell is a FIFO data store
	data   interface{}
	lock   sync.Mutex
	parent *streamTable
}

func (st *streamTable) push() { // push function will be called when entrycount == shape

	switch location := st.destination.(type) {
	case *table:
		buf := []interface{}{}
		for _, v := range st.cells {
			buf = append(buf, v.data)
		}
		location.Insert(buf)
	}
	st.entrycount = 0
	for _, v := range st.cells {
		v.lock.Unlock()
	}
}

func (cell *streamcell) push(data interface{}) {
	cell.lock.Lock() // wait until table is pushed
	cell.data = data
	cell.parent.entrycount++
	if cell.parent.shape == cell.parent.entrycount {
		cell.parent.push()
	}
}
