package main

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

func Sql_Handler(commands []string) (*string, map[string]interface{}, error) {

	/* syntax check
	// data key check
	// seperate process columns
	// register clauses
	// register functions
	// interprete to native nosql and golang func
	// register task flow nodes
	// execute nodes in order
	// todo: optimization by monitering each node

	type node struct { // use a node struct to clearify task order
		name string
		next *node
	}
	*/

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	length := len(commands)

	//sqlcode := strings.Join(commands, " ")

	//sqlcodeLines := strings.Split(sqlcode, ";")

	switch commands[0] {

	case "SELECT":

		/*
			steps:
			-check if table exist
			-check if table is a map
			-seperate each column statment
			-seperate clauses e.g. WHERE, ORDER BY
			-check if all columns has to be DISTINCT
			-seperate alias name
			-seperate and register functions

		*/

		//resaultmap := map[string]interface{}{}

		//tablename := commands[stringSliceIndex(commands, "FROM")+1]

	case "UPDATE":

	case "DELETE":

	case "INSERT":

	case "WHERE":

	case "CREATE":
		switch commands[1] {

		case "TABLE":
			numreg, err := regexp.Compile("[^a-zA-Z]+") // no number included
			if err != nil {
				log.Println(err)
				return nil, nil, err
			}

			if commands[3] == "AS" { // syntax: CREATE TABLE name AS SELECT column, column1 FROM table

				if commands[5] == "*" && commands[4] == "SELECT" {

					return nil, nil, nil
				}

				table := commands[length-1]

				if v, ok := shardData[table]; ok {

					if maps, ok := v.(map[string]interface{}); ok {

						columns := commands[4 : length-3]

						newmap := map[string]interface{}{}

						for i, v := range columns {
							columns[i] = strings.Replace(v, ",", "", -1)

							if a, ok := maps[columns[i]]; ok {

								if _, ok := table_type_map[table]; !(ok) { // check if is table
									return nil, nil, errors.New("table not exist")
								}
								switch table_type_map[table][columns[i]] {
								case "[]string":
									newmap[columns[i]] = a.([]string)
									table_type_map[commands[2]][columns[i]] = "[]string"

								case "[]int":
									newmap[columns[i]] = a.([]int)
									table_type_map[commands[2]][columns[i]] = "[]int"

								case "[]bool":
									newmap[columns[i]] = a.([]bool)
									table_type_map[commands[2]][columns[i]] = "[]bool"

								case "[]float":
									newmap[columns[i]] = a.([]float64)
									table_type_map[commands[2]][columns[i]] = "[]float"

								case "[][]byte":
									newmap[columns[i]] = a.([][]byte)
									table_type_map[commands[2]][columns[i]] = "[][]byte"

								default:
									return nil, nil, errors.New("invalid data type")
								}

							}

						}
						shardData[commands[2]] = newmap
					}
				}

			} else { // syntax: CREATE TABLE name (column int, column1 varchar)
				shardData[commands[2]] = map[string]interface{}{}
				columns := []string{}
				for i, v := range commands[3 : len(commands)-1] {
					columns[i] = reg.ReplaceAllString(v, "")
				}
				for i, v := range columns {
					switch numreg.ReplaceAllString(columns[i+1], "") {
					case "int":
						shardData[commands[2]].(map[string]interface{})[v] = []int{}

					case "varchar":
						shardData[commands[2]].(map[string]interface{})[v] = []string{}

					case "char":
						shardData[commands[2]].(map[string]interface{})[v] = []string{}

					case "varbinary":
						shardData[commands[2]].(map[string]interface{})[v] = [][]byte{}

					case "binary":
						shardData[commands[2]].(map[string]interface{})[v] = [][]byte{}

					case "bytes":
						shardData[commands[2]].(map[string]interface{})[v] = [][]byte{}

					case "float":
						shardData[commands[2]].(map[string]interface{})[v] = []float64{}

					case "bool":
						shardData[commands[2]].(map[string]interface{})[v] = []bool{}
					}
				}

			}

		}
	default:
		return nil, nil, errors.New("command not found")

	}
	return nil, nil, nil
}

func sql_func_register(function string, register map[int]string, column_index int) *string {
	funcs := []string{"COUNT", "MIN", "MAX", "AVG", "SUM", "SQRT", "RAND"}

	if contains(funcs, function) {
		register[column_index] = function

		return &function

	} else {
		return nil
	}
}
