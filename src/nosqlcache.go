/*
*************************************************************************
**This module records the hash of no sql command. When a same command is used frequently,
**this module will store the procedure into cache, looping through a switch will be unecessary
**the first time procedure stored into register, command will run as usual, only will the
**procedure be executed through cache next time.
 */
package main

type nosql_cache struct {
	hash      []byte
	command   string // Set, Get, Update, Clone, Move
	location  map[string]interface{}
	location1 map[string]interface{}
	key       string
	real_data interface{}
}

func (p *nosql_cache) execute() (string, error) {
	switch p.command {
	case "Set":
		cache_Set(p.location, p.key, p.real_data)
	}
	return "sucess", nil
}

func cache_Set(loc map[string]interface{}, key string, value interface{}) {
	shardData_lock.Lock()
	loc[key] = value
	shardData_lock.Unlock()
}
