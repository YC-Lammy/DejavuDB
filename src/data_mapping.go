package main

import (
	"errors"
	"net"
	"strings"
)

type data_map_key struct {
	subkeys    map[string]*data_map_key // keyname: key
	data_group *data_replicates_group   //data_management.go
	permission int                      // 3 digit permission number user, group, others
	owner      int
	group      int
}

var data_map = map[string]interface{}{} // key.key.key = []string, only map[string]interface{} and []string

var data_map1 = map[string]*data_map_key{}

func getKey(location string) (*data_map_key, error) {
	keys := strings.Split(location, ".")
	var pointer map[string]*data_map_key = data_map1
	var mapkey *data_map_key

	for _, v := range keys { // lopthrough every key until a data group is found
		if key, ok := pointer[v]; ok {
			pointer = key.subkeys // redirect pointer to subkey
			mapkey = key
		} else {
			return nil, errors.New("invalid data key, key " + v + " undefined")
		}
	}
	return mapkey, nil
}

func chmod(command string) error {
	return nil
}

func chown(command string) error {
	return nil
}

func chgrp(command string) error {
	return nil
}

func getShardConn(location string) ([]net.Conn, error) {

	keys := strings.Split(location, ".")
	var pointer map[string]*data_map_key = data_map1

	for _, v := range keys { // lopthrough every key until a data group is found
		if key, ok := pointer[v]; ok {
			pointer = key.subkeys // redirect pointer to subkey
			if g := key.data_group; g != nil {
				return g.connections, nil
			}
		} else {
			return nil, errors.New("invalid data key, key " + v + " undefined")
		}
	}
	// no groups registered in key range, function not returned
	switch len(pointer) {

	case 0:
		return nil, errors.New("key initiated without any data") // no subkey and data, empty ghost key

	default:
		conns := []net.Conn{}
		keybuffer := []map[string]*data_map_key{pointer}

		// loop through every subkeys and find all connections until no subkey is found
		for i := 0; i < len(keybuffer); i++ {
			for _, v := range keybuffer[i] {
				if v.data_group != nil {
					conns = append(conns, v.data_group.connections...)
				} else if len(v.subkeys) > 0 {
					keybuffer = append(keybuffer, v.subkeys)

				}
			}
		}
		return conns, nil
	}
}

func make_map_key(location string, data_group *data_replicates_group, permission, owner, group int) error {
	keys := strings.Split(location, ".")
	pointer := data_map1
	for i, v := range keys { // delete all existed keys from list and return a endpoint
		if _, ok := pointer[v]; ok {
			pointer = pointer[v].subkeys
			keys = keys[i+1:] // remove key
		} else { // the key does not exist
			break
		}
	}
	if len(keys) == 0 {
		return errors.New("key already exist")
	}
	for _, v := range keys { // loop through non exist keys
		pointer[v] = &data_map_key{subkeys: map[string]*data_map_key{}, data_group: data_group, permission: permission, owner: owner, group: group}
		pointer = pointer[v].subkeys
	}
	return nil
}
