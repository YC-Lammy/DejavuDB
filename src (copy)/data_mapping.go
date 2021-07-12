package main

import (
	"errors"
	"strings"
)

type data_map_key struct {
	subkeys    map[string]*data_map_key
	group      *data_replicates_group
	permission int
}

var data_map = map[string]interface{}{} // key.key.key = []string, only map[string]interface{} and []string

func getShardMac(location string) ([]string, error) { // get the mac addr of the shard that saves the data

	keys := strings.Split(location, ".")
	var pointer map[string]interface{}

	if v, ok := data_map[keys[0]]; ok {
		if i, ok := v.([]string); ok {
			return i, nil
		}
		if i, ok := v.(map[string]interface{}); ok {
			pointer = i
		} else {
			return nil, errors.New("type not match")
		}
	}

	for _, key := range keys[1:] {
		if v, ok := pointer[key]; ok {
			switch v := v.(type) {
			case map[string]interface{}:
				pointer = v

			case []string:
				return v, nil

			}

		}
	}
	buffer := []map[string]interface{}{}
	macs := []string{}
	// mac not found, find every mac under the pointer instead
	for _, v := range pointer {
		switch v := v.(type) {
		case []string:
			macs = append(macs, v...)
		case map[string]interface{}:
			buffer = append(buffer, v)
		default:
			return nil, errors.New("invalid type")

		}
	}
	if len(buffer) > 0 {
		for len(buffer) > 0 {
			for i, v := range buffer {
				buffer = append(buffer[:i], buffer[i+1:]...)
				for _, v := range v {
					switch v := v.(type) {
					case []string:
						macs = append(macs, v...)
					case map[string]interface{}:
						buffer = append(buffer, v)
					default:
						return nil, errors.New("invalid type")
					}
				}
			}
		}
	}
	return macs, nil
}
