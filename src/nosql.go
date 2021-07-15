package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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
					pointer = a.(map[string]interface{})
				} else {
					pointer[v] = map[string]interface{}{} // create a new key
					pointer = pointer[v].(map[string]interface{})
				}
			}
		}

		switch commands[3] {

		case "int":
			v, err := strconv.ParseInt(commands[2], 10, 64)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = int(v)

		case "float64", "float":
			v, err := strconv.ParseFloat(commands[2], 64)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

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
					pointer = a.(map[string]interface{})
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
			a := strconv.FormatInt(int64(v), 10)
			return &a, nil

		case float64:
			a := strconv.FormatFloat(v, 'g', -1, 64)
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
					pointer = a.(map[string]interface{})
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

		case float64:
			v, err := strconv.ParseFloat(commands[2], 64)
			if err != nil {
				return nil, err
			}
			pointer[keys[len(keys)-1]] = v

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
					pointer = a.(map[string]interface{})
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
					pointer = a.(map[string]interface{})
				} else {
					return nil, errors.New("key not exist")
				}

			}
		}
		pointer1 := shardData

		keys1 := strings.Split(commands[1], ".")

		if len(keys1) == 1 {
			if _, ok := pointer[keys1[0]]; !ok {
				pointer1[keys1[0]] = map[string]interface{}{} // create a new key
				pointer1 = pointer1[keys1[0]].(map[string]interface{})
			}
		} else {

			for _, v := range keys1[:len(keys1)-1] {
				if a, ok := pointer1[v]; ok {
					pointer1 = a.(map[string]interface{})
				} else {
					pointer1[v] = map[string]interface{}{} // create a new key
					pointer1 = pointer1[v].(map[string]interface{})
				}
			}
		}

		switch v := pointer[keys[len(keys)-1]].(type) {

		case string, int, float64, bool, []byte, []string, []int, [][]byte, []float64, []bool, map[string]interface{}:
			pointer1[keys1[len(keys1)-1]] = v
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
	}
	return nil, errors.New("command not found")
}
