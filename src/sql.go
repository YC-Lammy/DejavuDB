package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func Sql_Handler(commands []string) (*string, map[string]interface{}, error) { // asume syntax checked

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
			shardData[tabelename] = tmp_table
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
					if strings.Count(value, "(") <= 1 {
						function_split := strings.Split(value, "(")

						haveFunc := sql_func_register(function_split[0], apply_function_on_column, column_index)
					} else {

					}

					if haveFunc != nil {
						value = strings.Replace(function_split[1], ")", "", 1)

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
									return nil, errors.New("table not exist")
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

func sql_MAX(array interface{}) ([]float64, error) {
	result := []float64{}
	var tmp float64
	switch v := array.(type) {
	case []int:

		tmp = math.Max(float64(v[0]), float64(v[1]))
		for i := 0; i < len(v); i++ {
			tmp = math.Max(tmp, float64(v[i]))
		}
		result = append(result, tmp)
		return result, nil

	case []float64:

		tmp = math.Max(float64(v[0]), float64(v[1]))
		for i := 0; i < len(v); i++ {
			tmp = math.Max(tmp, float64(v[i]))
		}
		result = append(result, tmp)
		return result, nil

	default:
		return nil, errors.New("sql: MAX function expected array")
	}
}

func sql_MIN(array interface{}) ([]float64, error) {
	result := []float64{}
	var tmp float64
	switch v := array.(type) {
	case []int:

		tmp = math.Min(float64(v[0]), float64(v[1]))
		for i := 0; i < len(v); i++ {
			tmp = math.Min(tmp, float64(v[i]))
		}
		result = append(result, tmp)
		return result, nil

	case []float64:

		tmp = math.Min(float64(v[0]), float64(v[1]))
		for i := 0; i < len(v); i++ {
			tmp = math.Min(tmp, float64(v[i]))
		}
		result = append(result, tmp)
		return result, nil

	default:
		return nil, errors.New("sql: MIN function expected array")
	}
}

func sql_AVG(array interface{}) ([]float64, error) {
	result := []float64{}
	var tmp float64 = 0
	switch v := array.(type) {
	case []int:
		for _, value := range v {
			tmp += float64(value)
		}

		tmp = tmp / float64(len(v))

	case []float64:
		for _, value := range v {
			tmp += value
		}
		tmp = tmp / float64(len(v))
	default:
		return nil, errors.New("sql: AVG function expected array")
	}

	result = append(result, tmp)

	return result, nil
}

func sql_SUM(array interface{}) ([]float64, error) {
	result := []float64{}
	var tmp float64 = 0
	switch v := array.(type) {
	case []int:
		for _, value := range v {
			tmp += float64(value)
		}

	case []float64:
		for _, value := range v {
			tmp += value
		}

	default:
		return nil, errors.New("sql: SUM function expected array")
	}

	result = append(result, tmp)

	return result, nil
}

// Start of sql string functions

func sql_CHAR_LENGTH(array []string) ([]int, error) {
	result := []int{}
	for _, v := range array {
		result = append(result, len(v))
	}

	return result, nil
}

func sql_CHARACTER_LENGTH(array []string) ([]int, error) {
	return sql_CHAR_LENGTH(array)
}

func sql_LCASE(array []string) ([]string, error) {
	result := []string{}
	for _, v := range array {
		result = append(result, strings.ToLower(v))

	}
	return result, nil
}
func sql_LOWER(array []string) ([]string, error) {
	return sql_LCASE(array)
}

func sql_UCASE(array []string) ([]string, error) {
	result := []string{}
	for _, v := range array {
		result = append(result, strings.ToUpper(v))

	}
	return result, nil
}

func sql_UPPER(array []string) ([]string, error) {
	return sql_UCASE(array)
}

func sql_LENGTH(array []string) ([]int, error) {
	result := []int{}
	for _, v := range array {
		result = append(result, len(v))
	}
	return result, nil
}

func sql_REVERSE(array []string) ([]string, error) {
	result := []string{}
	for _, v := range array {
		var str = ""
		for _, v := range strings.Split(v, "") {
			str = v + str
		}
		result = append(result, str)
	}
	return result, nil
}

/*
END of sql string functions
*/
func sql_CONCAT(columns []interface{}) ([]string, error) {
	result := []string{}
	var num int = 0
	switch v := columns[0].(type) {
	case []int:
		num = len(v)
	case []string:
		num = len(v)
	case []bool:
		num = len(v)
	case [][]byte:
		num = len(v)
	case []float64:
		num = len(v)
	}
	for i := 0; i < num; i++ {
		result = append(result, "")
	}
	for _, y := range columns {
		switch v := y.(type) {
		case []int:
			for i, a := range v {
				result[i] += strconv.FormatInt(int64(a), 10)
			}

		case []string:
			for i, a := range v {
				result[i] += a
			}
		case []bool:
			for i, a := range v {
				result[i] += strconv.FormatBool(a)
			}
		case [][]byte:
			for i, a := range v {
				result[i] += string(a)
			}
		case []float64:
			for i, a := range v {
				result[i] += strconv.FormatFloat(a, 'g', -1, 64)
			}
		}
	}
	return result, nil

}
func sql_CONCAT_WS(columns []interface{}, ws string) ([]string, error) {

	result := []string{}
	var num int = 0
	switch v := columns[0].(type) {
	case []int:
		num = len(v)
	case []string:
		num = len(v)
	case []bool:
		num = len(v)
	case [][]byte:
		num = len(v)
	case []float64:
		num = len(v)
	}
	for i := 0; i < num; i++ {
		result = append(result, "")
	}
	for _, y := range columns {
		switch v := y.(type) {
		case []int:
			for i, a := range v {
				result[i] += strconv.FormatInt(int64(a), 10)
				if i < num-1 {
					result[i] += ws
				}
			}

		case []string:
			for i, a := range v {
				result[i] += a
				if i < num-1 {
					result[i] += ws
				}
			}
		case []bool:
			for i, a := range v {
				result[i] += strconv.FormatBool(a)
				if i < num-1 {
					result[i] += ws
				}
			}
		case [][]byte:
			for i, a := range v {
				result[i] += string(a)
				if i < num-1 {
					result[i] += ws
				}

			}
		case []float64:
			for i, a := range v {
				result[i] += strconv.FormatFloat(a, 'g', -1, 64)
				if i < num-1 {
					result[i] += ws
				}
			}
		}
	}
	return result, nil
}

// start of sql Numeric functions

func sql_ABS(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Abs(v))

	}
	return result, nil
}

func sql_ACOS(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Acos(v))
	}
	return result, nil
}

func sql_ACOSH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Acosh(v))
	}
	return result, nil
}

func sql_ASIN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Asin(v))

	}
	return result, nil
}

func sql_ASINH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Asinh(v))

	}
	return result, nil
}
func sql_ATAN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Atan(v))

	}
	return result, nil
}
func sql_ATANH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Atanh(v))

	}
	return result, nil
}

func sql_CEIL(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Ceil(v))

	}
	return result, nil
}

func sql_COS(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Cos(v))

	}
	return result, nil
}

func sql_COSH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Cosh(v))

	}
	return result, nil
}

func sql_CBRT(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Cbrt(v))

	}
	return result, nil
}

func sql_EXP(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Exp(v))

	}
	return result, nil
}

func sql_EXP2(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Exp2(v))

	}
	return result, nil
}

func sql_EXMP1(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Expm1(v))

	}
	return result, nil
}

func sql_FLOOR(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Floor(v))

	}
	return result, nil
}

func sql_LN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Log(v))

	}
	return result, nil
}

func sql_LOG(array []float64) ([]float64, error) {
	return sql_LN(array)
}

func sql_LOG10(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Log10(v))

	}
	return result, nil
}

func sql_LOG2(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Log2(v))

	}
	return result, nil
}

func sql_SIGNBIT(array []float64) ([]bool, error) {
	result := []bool{}
	for _, v := range array {
		result = append(result, math.Signbit(v))

	}
	return result, nil
}

func sql_SIN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Sin(v))

	}
	return result, nil
}

func sql_SINH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Sinh(v))

	}
	return result, nil
}

func sql_SQRT(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Sqrt(v))

	}
	return result, nil
}

func sql_TAN(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Tan(v))

	}
	return result, nil
}

func sql_TANH(array []float64) ([]float64, error) {
	result := []float64{}
	for _, v := range array {
		result = append(result, math.Tanh(v))

	}
	return result, nil
}
