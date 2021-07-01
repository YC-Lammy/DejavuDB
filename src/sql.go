package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"
)

func Sql_Handler(commands []string) (*string, map[string]interface{}, error) { // asume syntax checked

	type node struct { // use a node struct to clearify task order
		name string
		next *node
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	length := len(commands)

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

		resaultmap := map[string]interface{}{}

		tablename := commands[stringSliceIndex(commands, "FROM")+1]

		if tablename == "(SELECT" {
			n := 5
			b := make([]byte, n)
			if _, err := rand.Read(b); err != nil {
				return nil, nil, err
			}
			tablename = string(b)
			//                               split to []string           join to split ")"                                             remove "("
			_, tmp_table, err := Sql_Handler(strings.Split(strings.Split(strings.Join(commands[stringSliceIndex(commands, "FROM")+1:], " ")[1:], ")")[0], " "))
			if err != nil {
				return nil, nil, err
			}
			shardData[tablename] = tmp_table
			defer delete(shardData, tablename)
		}

		if v, ok := shardData[tablename]; ok {

			if table, ok := v.(map[string]interface{}); ok {

				operation_columns := strings.Split(strings.Join(commands[1:stringSliceIndex(commands, "FROM")-1], ""), ",")

				var step_map = map[string]int{}

				var distinct_all_columns = false

				var apply_function_on_column = map[int]string{}

				var distinct_on_single_column = map[int]string{}

				var headings = []string{}

				for _, key := range []string{"ORDER", "WHERE", "INNER", "LEFT", "RIGHT", "FULL", "GROUP", "UNION", "ORDER", "HAVING"} {
					if i := stringSliceIndex(commands, key); i != -1 {
						step_map[key] = i
					}
				}

				for column_index, value := range operation_columns {

					if value[0:8] == "DISTINCT" && column_index == 0 { // syntax: SELECT DISTINCT coumn FROM name

						distinct_all_columns = true
						value = strings.Replace(value, "DISTINCT", "", -1)

					}

					s := strings.Split(value, "AS")
					heading := ""

					if len(s) > 2 {
						return nil, nil, errors.New("sql: syntax error")
					}

					if len(s) == 2 {
						headings = append(headings, s[1])
						heading = s[1]
					}

					if strings.Count(value, "(") == 1 {
						function_split := strings.Split(value, "(")

						haveFunc := sql_func_register(function_split[0], apply_function_on_column, column_index)

						if haveFunc != nil {
							value = strings.Replace(function_split[1], ")", "", 1)
						} else {
							return nil, nil, errors.New("sql: function not found")
						}
					} else if strings.Count(value, "(") == 0 {

					} else {

						function_split := strings.Split(value, "(")

						for _, v := range function_split[:len(function_split)-1] {

							sql_func_register(v, apply_function_on_column, column_index)
						}

						value = strings.Replace(function_split[len(function_split)-1], ")", "", -1)
					}

					if value[0] == '*' {
						resaultmap = table
						break
					}

					if a, ok := table[value]; ok {

						switch v := a.(type) {

						case []string:
							resaultmap[value] = v
						case []int:
							resaultmap[value] = v
						case []bool:
							resaultmap[value] = v
						case []float64:
							resaultmap[value] = v
						case [][]byte:
							resaultmap[value] = v
						default:
							return nil, nil, errors.New("sql: invalid data type")

						}
					}

				} // end loop

				a, err := json.Marshal(resaultmap)

				if err != nil {
					return nil, nil, err
				}

				result := string(a)

				return &result, resaultmap, nil

			} else {
				return nil, nil, errors.New("invalid type of interface")
			} // end if

		} else {
			return nil, nil, errors.New("table not exist")
		}
		//end if

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
