package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/DmitriyVTitov/size"
)

func ShardHandler(conn net.Conn, message string) {
	defer wg.Done()

	commands := strings.Split(message, " ")
	var resault string
	switch commands[0] {

	case "groupadd":
		groupadd(conn, commands) // groupadd [option] groupName

	case "useradd":
		useradd(conn, commands) // useradd [option] userName

	case "size":
		resault = string(getShardSize())

	default:
		fmt.Fprintln(conn, "Invalid")
		return
	}
	fmt.Fprintln(conn, resault)
}

func groupadd(conn net.Conn, commands []string) { // option -g specified group id, -r system group

	_, exist := shardData[commands[len(commands)-1]]

	var id int64 = 1001
	var err error

	if exist { // group name exist
		fmt.Fprintln(conn, "Invalid")
		return
	}
	if len(commands) > 2 {
		for i := 0; i < len(commands); i++ {
			if commands[i] == "-g" {
				id, err = strconv.ParseInt(commands[i+1], 10, 64)
				CheckErr(err)
			}
			if commands[i] == "-r" {
				id = 50
			}
		}
	}

	shardData[commands[len(commands)-1]] = map[string]interface{}{"id": id}

}

func useradd(conn net.Conn, commands []string) { // option -u specified user id, -G add to group

	var id int64 = 1001
	var err error
	var group string = "user"
	var username string = commands[len(commands)-1]

	if len(commands) > 2 {
		for i := 0; i < len(commands); i++ {
			if commands[i] == "-u" {
				id, err = strconv.ParseInt(commands[i+1], 10, 64)
				CheckErr(err)
			}
			if commands[i] == "-G" {
				group = commands[i+1]
			}
		}
	}
	_, exist := shardData[group].(map[string]interface{})[username]

	if exist { // username name exist
		fmt.Fprintln(conn, "Invalid")
		return
	}

	shardData[group].(map[string]interface{})[username] = map[string]interface{}{"id": id, "issue_date": time.Now().String()}
}

func getShardSize() int {
	return size.Of(shardData)
}
