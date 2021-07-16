package main

type table struct {
	name          string
	column_dtypes map[string]string
	columns       map[string]*column
	permission    [3]int8 // 0-7, 3 digit permission number owner, group, others e.g. 770
	owner         int     // owner id
	group         int     // group id
}

type column struct {
	name     string
	datatype string
	data     interface{}
}

func (tb table) WriteRow(values ...interface{}) {

}
