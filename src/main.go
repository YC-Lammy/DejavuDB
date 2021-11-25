package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"sync"

	"net/http"
	_ "net/http/pprof"

	"dejavuDB/src/config"
	"dejavuDB/src/lazy"

	"dejavuDB/src/standalone"
	"dejavuDB/src/static"
)

var mycfg = []byte{}

var MAC_Address string = lazy.Get_first_mac_addr() // get.go

var sql_file string = ""

var wg sync.WaitGroup // working group

func main() {

	if config.Debug {
		go http.ListenAndServe(":6060", nil)
	}

	for _, v := range os.Args {
		if v == "help" {
			fmt.Println(static.Manual)
			return
		}
	}

	gob.Register(map[string]interface{}{})

	//fmt.Println("enter your password:")

	//fmt.Scanln(&password)

	fmt.Println("save to disk: ", config.Save_disk)

	home_dir, _ := os.UserHomeDir()

	os.Chdir(home_dir)

	os.Mkdir("dejavuDB", os.ModePerm)

	os.Chdir("dejavuDB")

	os.Mkdir("log", os.ModePerm)
	os.Mkdir("database", os.ModePerm)
	os.Mkdir("ML", os.ModePerm)
	os.Mkdir("addon", os.ModePerm)

	os.Chdir(path.Join(home_dir, "dejavuDB"))

	os.Chdir("database")

	os.Chdir(path.Join(home_dir, "dejavuDB"))

	if config.Password != "a empty password" {
		for {
			if len(config.Password) != 16 && len(config.Password) != 24 && len(config.Password) != 32 {
				fmt.Println("password must be length of 16, 24 or 32")
				fmt.Scanln(&config.Password)
			} else {
				break
			}
		}

	}

	fmt.Println("role: " + config.Role + " listener ip: " + config.Leader_addr)
	switch config.Role {

	case "router":

	case "shard":

	case "client":

	case "standalone":
		go standalone.Start()
		wg.Add(1)

	case "log":

	default:
		panic("Specified Role Invalid")

	}

	wg.Wait() // wait until all worker end

}
