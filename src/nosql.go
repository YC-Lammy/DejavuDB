package main

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/DmitriyVTitov/size"
)

func Nosql_Handler(commands []string) (*string, error) {

	sucess := "sucess"

	switch commands[0] {

	case "Set": // syntax: "Set name.name1.name2 value type" e.g. "Set User.John.id 23740 int"

		if len(commands) < 4 {
			return nil, errors.New("3 arguements required")
		}

		pointer := shardData

		keys := strings.Split(commands[1], ".")

		if len(keys) != 1 {
			for _, v := range keys[:len(keys)-1] {
				if a, ok := pointer[v]; ok {
					if b, ok := a.(map[string]interface{}); ok {
						pointer = b
					} else {
						pointer[v] = map[string]interface{}{} // overwrite key store value
						pointer = pointer[v].(map[string]interface{})
					}
				} else {
					pointer[v] = map[string]interface{}{} // create a new key
					pointer = pointer[v].(map[string]interface{})
				}
			}
		}

		switch commands[3] {

		case "int", "int64":
			v, err := strconv.Atoi(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = int(v)

		case "int128":
			v, err := strToInt128(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

		case "int256":
			v, err := strToInt256(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

		case "float64", "float", "ft", "ft64":
			v, err := strconv.ParseFloat(commands[2], 64)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

		case "bigfloat", "big_float", "bf":
			f, _, err := new(big.Float).Parse(commands[2], 10)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = f

		case "longdouble", "long_double", "float128", "ft128": // store value in c.longdouble
			f, err := strToFloat128(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = f

		case "string", "str":
			pointer[keys[len(keys)-1]] = commands[2]

		case "bool":
			v, err := strconv.ParseBool(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

		case "[]byte", "bytes":
			pointer[keys[len(keys)-1]] = []byte(commands[2])

		case "[]string", "str_arr":

			str := commands[2]
			result := []string{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			result = append(result, a...)

			pointer[keys[len(keys)-1]] = result

		case "[]int", "int_arr":
			str := commands[2]
			result := []int{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				b, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, err
				}
				result = append(result, int(b))
			}

			pointer[keys[len(keys)-1]] = result

		case "[]float64", "[]float", "float_arr":
			str := commands[2]
			result := []float64{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				b, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, err
				}
				result = append(result, b)
			}

			pointer[keys[len(keys)-1]] = result

		case "[][]byte", "bytes_arr":
			str := commands[2]
			result := [][]byte{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				result = append(result, []byte(v))
			}

			pointer[keys[len(keys)-1]] = result

		case "[]bool", "bool_arr":
			str := commands[2]
			result := []bool{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, " ", "", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, err
				}
				result = append(result, b)
			}

			pointer[keys[len(keys)-1]] = result

		case "json":
			var result = map[string]interface{}{}
			err := json.Unmarshal([]byte(commands[2]), &result)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = result

		default:
			return nil, errors.New("invalid type " + commands[3])

		}

		return &sucess, nil

	case "Get": // syntax: "Get name.name1.name2" e.g. "Get User.John.id" -> 23740

		if len(commands) < 2 {
			return nil, errors.New("1 argument required")
		} else if len(commands) > 2 {
			return nil, errors.New("only 1 argument required")
		}

		pointer := shardData

		keys := strings.Split(commands[1], ".")

		if len(keys) == 1 {
			if _, ok := pointer[keys[0]]; !ok {
				return nil, errors.New("key not exist")
			}
		} else {
			for _, v := range keys[:len(keys)-1] {
				if a, ok := pointer[v]; ok {
					if b, ok := a.(map[string]interface{}); ok {
						pointer = b
					} else {
						return nil, errors.New("key not exist")
					}
				} else {
					return nil, errors.New("key not exist")
				}
			}
		}

		switch v := pointer[keys[len(keys)-1]].(type) {

		case string:
			a := v
			return &a, nil

		case int:
			a := strconv.Itoa(v)
			return &a, nil

		case int8:
			a := strconv.Itoa(int(v))
			return &a, nil
		case int16:
			a := strconv.Itoa(int(v))
			return &a, nil
		case int32:
			a := strconv.Itoa(int(v))
			return &a, nil
		case int64:
			a := strconv.Itoa(int(v))
			return &a, nil

		case int128:
			a := v.String()
			return &a, nil

		case int256:
			a := v.String()
			return &a, nil

		case float64:
			a := fmt.Sprintf("%v", v)
			return &a, nil

		case float128:
			a := v.String()
			return &a, nil

		case *big.Float:
			a := v.Text('g', int(v.Prec()))
			return &a, nil

		case bool:
			a := strconv.FormatBool(v)
			return &a, nil

		case []byte:
			a := string(v)
			return &a, nil

		case []string:
			a := "[" + strings.Join(v, ",") + "]"
			return &a, nil

		case []int:
			b := []string{}
			for _, s := range v {
				b = append(b, strconv.FormatInt(int64(s), 10))
			}
			a := "[" + strings.Join(b, ",") + "]"
			return &a, nil

		case [][]byte:
			b := []string{}
			for _, s := range v {
				b = append(b, string(s))
			}
			a := "[" + strings.Join(b, ",") + "]"
			return &a, nil

		case []float64:
			b := []string{}
			for _, s := range v {
				b = append(b, strconv.FormatFloat(float64(s), 'g', -1, 64))
			}
			a := "[" + strings.Join(b, ",") + "]"
			return &a, nil

		case []bool:
			b := []string{}
			for _, s := range v {
				b = append(b, strconv.FormatBool(s))
			}
			a := "[" + strings.Join(b, ",") + "]"
			return &a, nil

		case map[string]interface{}:
			b, err := json.Marshal(v)
			a := string(b)
			if err != nil {
				return nil, err
			}
			return &a, nil
		}

	case "Update": // syntax: "Update name.name1.name2 value" e.g. "Update User.John.id 23740"

		pointer := shardData

		keys := strings.Split(commands[1], ".")

		if len(commands) < 3 {
			return nil, errors.New("2 arguements required")
		}

		if len(keys) == 1 {
			if _, ok := pointer[keys[0]]; !ok {
				return nil, errors.New("key " + keys[0] + " not exist")
			}
		} else {
			for _, v := range keys[:len(keys)-1] {
				if a, ok := pointer[v]; ok {
					if b, ok := a.(map[string]interface{}); ok {
						pointer = b
					} else {
						return nil, errors.New("key " + v + " is not a map")
					}
				} else {
					return nil, errors.New("key " + v + " not exist")
				}
			}

		}

		switch pointer[keys[len(keys)-1]].(type) {

		case int:
			v, err := strconv.ParseInt(commands[2], 10, 64)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = int(v)

		case int8:
			v, err := strconv.ParseInt(commands[2], 10, 8)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = int8(v)

		case int16:
			v, err := strconv.ParseInt(commands[2], 10, 16)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = int16(v)

		case int32:
			v, err := strconv.ParseInt(commands[2], 10, 32)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = int32(v)

		case int64:
			v, err := strconv.ParseInt(commands[2], 10, 64)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = int64(v)

		case int128:
			v, err := strToInt128(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v
		case int256:
			v, err := strToInt256(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

		case float32:
			v, err := strconv.ParseFloat(commands[2], 32)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = float32(v)

		case float64:
			v, err := strconv.ParseFloat(commands[2], 64)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

		case float128:
			v, err := strToFloat128(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

		case *big.Float:
			f, _, err := new(big.Float).Parse(commands[2], 10)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = f

		case string:
			pointer[keys[len(keys)-1]] = commands[2]

		case bool:
			v, err := strconv.ParseBool(commands[2])
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

		case []byte:
			pointer[keys[len(keys)-1]] = []byte(commands[2])

		case []string:

			str := commands[2]
			result := []string{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			result = append(result, a...)

			pointer[keys[len(keys)-1]] = result

		case []int:
			str := commands[2]
			result := []int{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				b, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, err
				}
				result = append(result, int(b))
			}

			pointer[keys[len(keys)-1]] = result

		case []float32:
			str := commands[2]
			result := []float32{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				b, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return nil, err
				}
				result = append(result, float32(b))
			}

			pointer[keys[len(keys)-1]] = result

		case []float64:
			str := commands[2]
			result := []float64{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				b, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, err
				}
				result = append(result, b)
			}

			pointer[keys[len(keys)-1]] = result

		case [][]byte:
			str := commands[2]
			result := [][]byte{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, ", ", ",", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				result = append(result, []byte(v))
			}

			pointer[keys[len(keys)-1]] = result

		case []bool:
			str := commands[2]
			result := []bool{}

			if commands[2][0] == '[' {
				str = commands[2][1 : len(commands[2])-1]
			}
			str = strings.Replace(str, " ", "", -1)

			a := strings.Split(str, ",")

			for _, v := range a {
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, err
				}
				result = append(result, b)
			}

			pointer[keys[len(keys)-1]] = result

		case map[string]interface{}:
			var result = map[string]interface{}{}
			err := json.Unmarshal([]byte(commands[2]), &result)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = result

		default:
			return nil, errors.New("no matching type")
		}
		return &sucess, nil

	case "Delete": // syntax: "Delete name.name1.name2" e.g. "Delete User.John.id"

		if len(commands) != 2 {
			return nil, errors.New("1 argument required")
		}

		pointer := shardData

		keys := strings.Split(commands[1], ".")

		if len(keys) == 1 {
			if _, ok := pointer[keys[0]]; !ok {
				return nil, errors.New("key " + keys[0] + " not exist")
			}
		} else {
			for _, v := range keys[:len(keys)-1] {
				if a, ok := pointer[v]; ok {
					if b, ok := a.(map[string]interface{}); ok {
						pointer = b
					} else {
						return nil, errors.New("key " + v + " is not a map")
					}
				} else {
					return nil, errors.New("key not exist")
				}

			}
		}

		delete(pointer, keys[len(keys)-1])
		return &sucess, nil

	case "Clone": // syntax: Clone target destination
		pointer := shardData

		keys := strings.Split(commands[1], ".")

		if len(keys) == 1 {
			if _, ok := pointer[keys[0]]; !ok {
				return nil, errors.New("key " + keys[0] + " not exist")
			}
		} else {

			for _, v := range keys[:len(keys)-1] {
				if a, ok := pointer[v]; ok {
					if b, ok := a.(map[string]interface{}); ok {
						pointer = b
					} else {
						return nil, errors.New("key " + v + " is not a map")
					}
				} else {
					return nil, errors.New("key not exist")
				}

			}
		}
		value := pointer[keys[len(keys)-1]]
		pointer1 := shardData

		keys1 := strings.Split(commands[2], ".")

		if len(keys1) != 1 {

			for _, v := range keys1[:len(keys1)-1] {
				if a, ok := pointer1[v]; ok {
					if b, ok := a.(map[string]interface{}); ok {
						pointer1 = b

					} else { // this key stores a value
						pointer1[v] = map[string]interface{}{} // overwrite the key
						pointer1 = pointer1[v].(map[string]interface{})
					}
				} else {
					pointer1[v] = map[string]interface{}{} // create a new key
					pointer1 = pointer1[v].(map[string]interface{})
				}
			}
		}

		switch v := value.(type) {

		case string, int, int8, int16, int32, int64, int128, int256, float32, float64, float128, *big.Float, bool, []byte, []string, []int, [][]byte, []float64, []bool, map[string]interface{}:
			pointer1[keys1[len(keys1)-1]] = v
		default:
			return nil, errors.New("type not match")
		}
		return &sucess, nil

	case "Move": // syntax: Move target destination
		if len(commands) < 3 {
			return nil, errors.New("2 arguements required")
		}
		_, err := Nosql_Handler([]string{"Clone", commands[1], commands[2]}) // Move can be done by two process
		if err != nil {
			return nil, err
		}
		_, err = Nosql_Handler([]string{"Delete", commands[1]})
		if err != nil {
			return nil, err
		}

		return &sucess, nil

	case "Sizeof", "SizeOf": // stntax: Sizeof location
		if len(commands) < 2 {
			return nil, errors.New("1 argument required")
		}

		v, err := getPointer(commands[1])
		if err != nil {
			return nil, err
		}
		a := strconv.Itoa(size.Of(v)) + " byte"
		return &a, nil

	case "Typeof", "TypeOf":
		if len(commands) < 2 {
			return nil, errors.New("1 argument required")
		}
		return getTypeByLocation(commands[1])

	}
	return nil, errors.New("command not found")
}

func getPointer(location string) (interface{}, error) {
	pointer := shardData

	keys := strings.Split(location, ".")

	if len(keys) == 1 {
		if _, ok := pointer[keys[0]]; !ok {
			return nil, errors.New("key " + keys[0] + " not exist")
		}
	} else {
		for _, v := range keys[:len(keys)-1] {
			if a, ok := pointer[v]; ok {
				if b, ok := a.(map[string]interface{}); ok {
					pointer = b
				} else {
					return nil, errors.New("key " + v + " is not a map")
				}
			} else {
				return nil, errors.New("key " + v + " not exist")
			}
		}

	}
	return pointer[keys[len(keys)-1]], nil
}

func getTypeByLocation(location string) (*string, error) {
	v, err := getPointer(location)
	if err != nil {
		return nil, err
	}
	a := fmt.Sprintf("%T", v)
	return &a, nil
}
