package main

import (
	"database/sql"
	"encoding/json"
)

type sql_dataset []map[string][]string

func (ds sql_dataset) Json() (string, error) {
	re := ""
	if len(ds) == 0 {
		return "", nil
	}
	for _, v := range ds {
		json, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		re += string(json) + "\n"
	}
	return re, nil
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
