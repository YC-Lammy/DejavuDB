package main

import (
	"net"
	"strings"
	"time"
)

// commands can only executed by a console client
// a full list of the commands
var list_of_commands = []string{"useradd", "groupadd", "Login", "Logout",
	"atop", "cat", "df", "dstat", "find", "free", "last", "netstat", "rm", "sort", "w", "top", "tar"}

var command_result string

func execute_command(conn net.Conn, message string) (string, error) {
	command_result = "" // reset result
	splited := strings.Split(message, " ")
	switch splited[0] {
	case "df":
		console_monitor(message, console_df_exec)
	case "free":
		console_monitor(message, console_free_exec)
	}
	return command_result, nil
}

// useradd in user.go
// groupadd in user.go
// Login in user.go
// Logout in user.go

func all_monitor_recieved() bool {
	return len(monitor_values) == router_connected+1
}

func console_monitor(message string, execfunc func()) {
	monitor_values = []map[string]string{}
	monitor_values = append(monitor_values, monitor())
	send_to_all_router([]byte("monitor"))
	send_to_all_shard([]byte("monitor"))
	waitUntil(all_monitor_recieved, execfunc, 10*time.Millisecond) // waitUtil lazy.go
}

func console_df_exec() { // exec of the df command
	str := "Role\tTotal\tUsed\tAvaliable\tUse%\n"
	for _, v := range monitor_values {

		str += strings.Join([]string{v["role"], v["disk_total"], v["disk_used"], v["disk_avaliable"], v["disk_load"]}, "\t") + "\n"

	}
	command_result = str

}

func console_free_exec() { // exec of the df command
	str := "Role\tTotal\tUsed\tAvaliable\tUse%\n"
	for _, v := range monitor_values {

		str += strings.Join([]string{v["role"], v["mem_total"], v["mem_used"], v["mem_avaliable"], v["mem_load"]}, "\t") + "\n"

	}
	command_result = str

}
