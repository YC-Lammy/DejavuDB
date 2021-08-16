package main

import (
	"errors"
	"sync"
)

/*
table type is a table that stores structured data.
table data type is not fixed, instead, it will dynamically change its data type according to input data
*/

// columns and rows in the same table shares the same set of cells in different direction
type column struct {
	name     string
	datatype byte    // each column can only have one data type
	data     []*cell // pointer to cell
}

type row struct {
	datatypes []byte
	data      []*cell
}

type cell struct {
	data     interface{}
	datatype byte
}

type table struct {
	name          string
	column_dtypes map[string]byte // map [column name] "data type"
	columns       []*column       // map [column name] pointer to column
	headers       []string
	rows          []*row
	column_count  int
	row_count     int

	permission [3]int8 // 0-7, 3 digit permission number owner, group, others e.g. 770
	owner      int     // owner id
	group      int     // group id

	insert_row_trigger_script string
}

func create_table(data string, name string, args ...interface{}) (*table, error) {
	if _, ok := test_shardData["tables"].key[name]; ok {
		return nil, errors.New("table name already used")
	}
	tab := table{name: name}
	new := &Node{data: &tab, lock: sync.Mutex{}}
	test_shardData["tables"].key[name] = new
	return &tab, nil
}

func (tb *table) Insert(Row []interface{}) error {
	if len(Row) != len(tb.columns) {
		return errors.New("column count mismatch")
	}
	newRow := row{data: []*cell{}}
	for _, v := range Row {
		newRow.data = append(newRow.data, &cell{data: v})
	}
	i := 0
	for _, v := range tb.columns {
		v.data = append(v.data, newRow.data[i])
		i++
	}
	return nil

}
