package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
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

//
var readOnlyPermission = [3]int8{7, 5, 5}
var publicPermission = [3]int8{7, 7, 7}
var groupPermission = [3]int8{7, 7, 0}
var ownerPermission = [3]int8{7, 0, 0}

//

var SQLDB = map[string]*sql_database{}
var defaultSQLDB *sql_database

var sqliteDB *sql.DB

type sql_database struct {
	name           string // name of database
	schemas        map[string]*schema
	permission     [3]int8
	owner          int
	group          int
	default_schema *schema // pointer to the schema, default :dbo
}

type schema struct {
	name       string
	tables     map[string]*table
	permission [3]int8
	owner      int
	group      int
}

type table struct {
	name          string
	column_dtypes map[string]string  // map [column name] "data type"
	columns       map[string]*column // map [column name] pointer to column
	permission    [3]int8            // 0-7, 3 digit permission number owner, group, others e.g. 770
	owner         int                // owner id
	group         int                // group id
}

func (tb *table) WriteRow(values ...interface{}) error {
	if len(values) != len(tb.columns) {
		return errors.New("column count doesn't match value count at row 1")
	}
	i := 0
	for _, c := range tb.columns {
		switch c.datatype {
		case "[]string":
			if v, ok := values[i].(string); ok {
				c.data = append(c.data.([]string), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]int":
			if v, ok := values[i].(int); ok {
				c.data = append(c.data.([]int), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]int8":
			if v, ok := values[i].(int8); ok {
				c.data = append(c.data.([]int8), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]int16":
			if v, ok := values[i].(int16); ok {
				c.data = append(c.data.([]int16), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]int32":
			if v, ok := values[i].(int32); ok {
				c.data = append(c.data.([]int32), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]int64":
			if v, ok := values[i].(int64); ok {
				c.data = append(c.data.([]int64), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]float32":
			if v, ok := values[i].(float32); ok {
				c.data = append(c.data.([]float32), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]float64":
			if v, ok := values[i].(float64); ok {
				c.data = append(c.data.([]float64), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]float128":
			if v, ok := values[i].(float128); ok {
				c.data = append(c.data.([]float128), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]bool":
			if v, ok := values[i].(bool); ok {
				c.data = append(c.data.([]bool), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[]byte":
			if v, ok := values[i].(byte); ok {
				c.data = append(c.data.([]byte), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		case "[][]byte":
			if v, ok := values[i].([]byte); ok {
				c.data = append(c.data.([][]byte), v)
			} else {
				return errors.New("there is a type mismatch at column " + strconv.Itoa(i+1))
			}
		}
		i++
	}
	return nil
}

type column struct {
	name     string
	datatype string
	data     interface{}
}

type views struct {
}

type procedures struct {
}

type sql_table map[string][]string // column name: [column data]

type sql_dataset []map[string][]string

func (ds sql_dataset) Json() (string, error) {
	re := map[string]map[string][]string{}
	for i, v := range ds {
		re["ResultSet"+strconv.Itoa(i+1)] = v
	}
	json, err := json.Marshal(re)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func read_SQL_Rows(rs *sql.Rows) (sql_dataset, error) {

	result := []map[string][]string{map[string][]string{}}

	names, err := rs.Columns()
	if err != nil {
		return nil, err
	}
	buffer := []interface{}{}
	for i := 0; i < len(names); i++ { // create exect number of pointer for scan
		a := ""
		buffer = append(buffer, &a) // create a new pointer
	}
	for rs.Next() { // loop to next row
		err = rs.Scan(buffer...) // scan takes execly numbers of columns
		if err != nil {
			return nil, err
		}
		for i, v := range names {
			if a, ok := result[0][v]; ok {
				result[0][v] = append(a, *buffer[i].(*string))
			} else {
				result[0][v] = []string{*buffer[i].(*string)}
			}
		}
	}
	s := 1
	for rs.NextResultSet() { // loop to next set
		result = append(result, map[string][]string{})
		names, err := rs.Columns()
		if err != nil {
			return nil, err
		}
		buffer := []interface{}{}
		for i := 0; i < len(names); i++ { // create exect number of pointer for scan
			a := ""
			buffer = append(buffer, &a) // create a new pointer
		}
		for rs.Next() {
			err = rs.Scan(buffer...) // scan takes execly numbers of columns
			if err != nil {
				return nil, err
			}
			for i, v := range names {
				if a, ok := result[s][v]; ok {
					result[s][v] = append(a, *buffer[i].(*string))
				} else {
					result[s][v] = []string{*buffer[i].(*string)}
				}
			}
		}
		s++
	}
	return sql_dataset(result), nil
}

func SQL_init() {

	sqliteDB, _ = sql.Open("sqlite3", ":memory:") // alternative sql_file

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
