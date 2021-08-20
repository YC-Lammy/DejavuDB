package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

/*
base data types:
	[]string
	[]int
	[]int8
	[]int16
	[]int32
	[]int64
	[]float32
	[]float64
	[]float128
	[]bool
	[]byte
	[][]byte
*/

var sqliteDB *sql.DB //, _ = sql.Open("sqlite3", ":memory:")

func SQL_init() {

	d, err := sql.Open("sqlite3", sql_file) // alternative sql_file
	if err != nil {
		panic(err)
	}
	sqliteDB = d // save pointer to global

}
