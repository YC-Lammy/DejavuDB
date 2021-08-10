package main

import (
	"flag"
)

type setting_fields struct {
	role               string
	router_addr        string
	password           string
	host               string
	port               string
	save_disk          bool
	debug              bool
	auto_shard         bool
	secure_connect     bool
	encrypt_disk       bool
	encrypt_bit_length bool
	enable_ML          bool
	javascript_timeout int // milli seconds
}

var Settings = setting_fields{}

func init_settings() {
	flag.StringVar(&Settings.role, "role", "full", "Specify role. Default is router, option: shard, client")
	flag.StringVar(&Settings.router_addr, "ip", "", "Specify router ip. Default is stand alone router")
	flag.StringVar(&Settings.password, "pwd", "a empty password", "Specify password. Default is empty")
	flag.StringVar(&Settings.host, "host", "localhost", "specify host addr")
	flag.StringVar(&Settings.port, "port", "8080", "specify hosting port")
	//flag.StringVar(&sql_file, "sqlfile", ":memory:", "specify sql file path")
	flag.BoolVar(&Settings.save_disk, "disk", false, "save copy to disk")
	flag.BoolVar(&Settings.debug, "debug", false, "developer debug option")
	flag.BoolVar(&Settings.enable_ML, "ML", false, "Enable built in machine learning service")
	flag.BoolVar(&Settings.secure_connect, "sc", false, "specify to use securite connection and the bit width")
	flag.IntVar(&Settings.javascript_timeout, "jstimeout", 500, "javascript vm timeout duration in milliseconds")
	flag.Parse()
}
