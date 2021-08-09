package main

import (
	"errors"
	"net"
	"strings"
	"time"
)

// commands can only executed by a console client
// a full list of the commands
var list_of_commands = []string{"Set", "Get", "Delete", "Update", "Clone", "Move", "Sizeof", "SizeOf", "Typeof", "TypeOf", "SQL",
	"useradd", "groupadd", "Login", "Logout",
	"atop", "cat", "cp", "chmod", "chown", "chgrp", "df", "dstat", "find",
	"free", "id", "last", "mv", "netstat", "rm", "sort", "w", "top", "tar"}

var command_result string

func execute_command(conn net.Conn, message string) (string, error) {
	command_result = "" // reset result
	splited := strings.Split(message, " ")
	var err error = nil
	switch splited[0] {
	case "Set", "Get", "Update", "Clone", "Move", "Delete", "Batch", "Sizeof", "SizeOf", "Typeof", "TypeOf", "SQL":
		router_apiHandler(conn, message)
		err = errors.New("do not send")
	case "useradd":
		err = useradd(message)
	case "groupadd":
		err = groupadd(message)
	case "Login":
	case "chmod":
		err = chmod(message)
	case "chown":
		err = chown(message)
	case "chgrp":
		err = chgrp(message)
	case "df":
		console_monitor(message, console_df_exec)
	case "dstat":
		console_monitor(message, console_dstat_exec)
	case "free":
		console_monitor(message, console_free_exec)
	case "id":
		err = userid(message)

	default:
		err = errors.New("command not found")
	}
	if err != nil {
		return "", err
	}
	return command_result, nil
}

// useradd in user.go
// groupadd in user.go
// Login in user.go
// Logout in user.go
// chmod data_mapping.go
// chown data_mapping.go
// chgrp data_mapping.go

func all_monitor_recieved() bool {
	return len(monitor_values) == router_connected+1
}

func console_monitor(message string, execfunc func()) {
	monitor_values = []map[string]string{}
	monitor_values = append(monitor_values, monitor())
	send_to_all_router([]byte("monitor"))
	send_to_all_shard([]byte("monitor"))
	waitUntil(all_monitor_recieved, execfunc, 1*time.Millisecond) // waitUtil lazy.go
}

func console_df_exec() { // exec of the df command
	str := "Disk\nRole\tTotal\tUsed\tAvaliable\tUse%\n"
	for _, v := range monitor_values {

		str += strings.Join([]string{v["role"], v["disk_total"], v["disk_used"], v["disk_avaliable"], v["disk_load"]}, "\t") + "\n"

	}
	command_result = str

}

func console_dstat_exec() { // exec of the df command
	str := "role,usage%,used,avaliable,use%,used,avaliable,use%;"
	for _, v := range monitor_values {

		str += strings.Join([]string{v["role"], v["cpu_load"] + "%", v["disk_used"] + "GB", v["disk_avaliable"] + "GB", v["disk_load"] + "%", v["mem_used"] + "GB", v["mem_avaliable"] + "GB", v["mem_load"] + "%"}, ",") + ";"

	}
	command_result = str

}

func console_free_exec() { // exec of the df command
	str := "memory\nRole\tTotal\tUsed\tAvaliable\tUse%\n"
	for _, v := range monitor_values {

		str += strings.Join([]string{v["role"], v["mem_total"], v["mem_used"], v["mem_avaliable"], v["mem_load"]}, "\t") + "\n"

	}
	command_result = str
}
