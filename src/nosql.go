package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func Nosql_Handler(commands []string) (*string, error) {
	switch commands[0] {
	case "Set": // syntax: "Set name.name1.name2 value type" e.g. "Set User.John.id 23740 int"

	case "Get": // syntax: "Get name.name1.name2" e.g. "Get User.John.id" -> 23740

		if len(commands) != 2 {
			return nil, errors.New("1 argument required")
		}

		pointer := shardData

		keys := strings.Split(commands[1], ".")

		for _, v := range keys[:len(keys)-1] {
			if a, ok := pointer[v]; ok {
				pointer = a.(map[string]interface{})
			} else {
				return nil, errors.New("key not exist")
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

		case []int:

		case [][]byte:

		case []float64:

		case []bool:

		case map[string]interface{}:
			b, err := json.Marshal(v)
			a := string(b)
			if err != nil {
				return nil, err
			}
			return &a, nil
		}

	case "Update": // syntax: "Update name.name1.name2 value type" e.g. "Update User.John.id 23740 int"

	case "Delete": // syntax: "Delete name.name1.name2" e.g. "Delete User.John.id"
		if len(commands) != 2 {
			return nil, errors.New("1 argument required")
		}

		pointer := shardData

		keys := strings.Split(commands[1], ".")

		for _, v := range keys[:len(keys)-1] {
			if a, ok := pointer[v]; ok {
				pointer = a.(map[string]interface{})
			} else {
				return nil, errors.New("key not exist")
			}

		}

		delete(pointer, keys[len(keys)-1])
	}
	return nil, nil
}
