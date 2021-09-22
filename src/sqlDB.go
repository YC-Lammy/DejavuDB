package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDB *sql.DB //, _ = sql.Open("sqlite3", ":memory:")

func SQL_init() {

	d, err := sql.Open("sqlite3", sql_file) // alternative sql_file
	if err != nil {
		panic(err)
	}
	sqliteDB = d // save pointer to global

}
