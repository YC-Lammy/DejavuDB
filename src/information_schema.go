package main

import "sync"

var information_schema = &Node{}

func init_information_schema() {
	information_schema.lock = sync.Mutex{}

	headers := []string{"PLUGIN_NAME", "PLUGIN_VERSION", "PLUGIN_STATUS", "PLUGIN_TYPE",
		"PLUGIN_TYPE_VERSION", "PLUGIN_LIBRARY", "PLUGIN_LIBRARY_VERSION", "PLUGIN_AUTHOR",
		"PLUGIN_DESCRIBSION", "PLUGIN_LICENSE", "LOAD_OPTION", "PLUGIN_MATURITY", "PLUGIN_AUTH_VERSION"}
	columns := []*column{}
	for _, v := range headers {
		columns = append(columns, &column{name: v, datatype: 0x01, data: []*cell{}})
	}
	rows := []*row{}
	for i := 0; i < len(headers); i++ {
		rows = append(rows, &row{})
	}
	all_plugnins_table := &table{name: "ALL_PLUGINS", headers: headers, columns: columns, rows: rows,
		permission: [3]int8{7, 4, 4}, owner: 100, group: 100}

	ALL_PLUGINS := Node{data: all_plugnins_table, lock: sync.Mutex{}}
	information_schema.key = map[string]*Node{"ALL_PLUGINS": &ALL_PLUGINS}
}
