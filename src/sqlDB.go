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

const (
	sql_str8    = 0x00 // byte
	sql_str     = 0x01 // []string
	sql_int     = 0x02 // []int32
	sql_int64   = 0x03 // []int64
	sql_float32 = 0x04 // []float32
	sql_float64 = 0x05 // []float64
	sql_bool    = 0x06 // []bool
)

//
var readOnlyPermission = [3]int8{7, 5, 5}
var publicPermission = [3]int8{7, 7, 7}
var groupPermission = [3]int8{7, 7, 0}
var ownerPermission = [3]int8{7, 0, 0}

//

var SQLDB = map[string]*sql_database{}
var defaultSQLDB *sql_database

var sqliteDB *sql.DB //, _ = sql.Open("sqlite3", ":memory:")

type sql_database struct {
	name           string // name of database
	schemas        map[string]*schema
	permission     [3]int8
	owner          int
	group          int
	default_schema *schema // pointer to the schema, default :dbo
}

type schema struct {
	name string

	tables        map[string]*table
	procedures    map[string]*procedure
	default_table *table

	permission [3]int8
	owner      int
	group      int
}

type views struct {
}

type procedure struct {
	hash   string
	opcode *Query
}

func SQL_init() {

	d, err := sql.Open("sqlite3", sql_file) // alternative sql_file
	if err != nil {
		panic(err)
	}
	sqliteDB = d // save pointer to global

	is := schema{}
	dbo := schema{permission: publicPermission}

	db1 := sql_database{
		name:           "Database_1",
		schemas:        map[string]*schema{"dbo": &dbo, "information_schema": &is},
		permission:     publicPermission,
		owner:          1000,
		group:          1000,
		default_schema: &dbo,
	}

	SQLDB["Database_1"] = &db1
	defaultSQLDB = &db1
}
