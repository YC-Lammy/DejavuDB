package main

type Tree struct {
	Next  *Tree
	Value int
}

var shardData = map[string]interface{}{}
var type_map = map[string]interface{}{}

var table_type_map = map[string]map[string]string{}
