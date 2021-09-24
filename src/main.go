package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	json "github.com/goccy/go-json"

	"src/lazy"

	"src/static"
)

var mycfg = []byte{}

var MAC_Address string = lazy.Get_first_mac_addr() // get.go

var sql_file string = ""

var wg sync.WaitGroup // working group

func main() {

	for i, v := range os.Args {
		switch v {
		case "help":
			fmt.Println(static.Manual)
			return
		case "install": // install add-on
			if _, err := os.Stat(os.Args[i+1]); os.IsNotExist(err) {
				// path/to/whatever does not exist
				panic(err)
			}
			return
		}
	}

	gob.Register(map[string]interface{}{})

	//fmt.Println("enter your password:")

	//fmt.Scanln(&password)

	fmt.Println("save to disk: ", config.Save_disk)

	os.Chdir(home_dir)

	os.Mkdir("dejavuDB", os.ModePerm)

	os.Chdir("dejavuDB")

	os.Mkdir("log", os.ModePerm)
	os.Mkdir("database", os.ModePerm)
	os.Mkdir("ML", os.ModePerm)
	os.Mkdir("addon", os.ModePerm)

	os.Chdir(path.Join(home_dir, "dejavuDB"))

	os.Chdir("database")

	os.Mkdir("tables", os.ModePerm)

	os.Chdir(path.Join(home_dir, "dejavuDB"))

	setupLog()

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
		start_router(config.Leader_addr)

	case "shard":
		start_shard(config.Leader_addr)

	case "client":
		start_client(config.Leader_addr)

	case "full":
		start_full(config.Leader_addr)

	case "log":
		start_log(config.Leader_addr)

	default:
		panic("Specified Role Invalid")

	}

	wg.Wait() // wait until all worker end

}

//////////////////////////////////////////////////////////////////////////////////////////

func start_router(dial_addr string) { // start as a router

	cfg := map[string]interface{}{"role": "router", "pass": config.Password, "mac": MAC_Address, "port": config.Host + ":" + config.Port}
	mycfg, _ = json.Marshal(cfg)

	//go process_timeout_checker()
	wg.Add(1)

	if dial_addr != "" {
		go start_listening() // router.go

		go dial_server(dial_addr, mycfg, RouterHandler, routerConfig) // network.go

		wg.Add(2)

	} else {

		fmt.Println("No ip specified, act as genesis router")

		go start_listening()

		wg.Add(1)
	}

}

/////////////////////////////////////////////////////////////////////////////////////

func start_shard(dial_addr string) { // start as a shard

	if dial_addr == "" {
		panic("must specific an address")
		return
	}
	cfg := map[string]interface{}{"role": "shard", "pass": config.Password, "mac": MAC_Address, "port": config.Host + ":" + config.Port}
	mycfg, _ = json.Marshal(cfg)

	go dial_server(dial_addr, mycfg, ShardHandler, shardConfig) // network.go

	go SQL_init()

	wg.Add(1)
}

////////////////////////////////////////////////////////////////////////////////////////

func start_client(dial_addr string) { // start as a client
	cfg := map[string]interface{}{"role": "client"}
	mycfg, _ = json.Marshal(cfg)
	go Client_dial(dial_addr, mycfg)

	wg.Add(1)
}

func start_log(dial_addr string) {
	if dial_addr == "" {
		panic("must specific an address")
		return
	}
	cfg := map[string]interface{}{"role": "log", "pass": config.Password, "port": config.Host + ":" + config.Port}
	mycfg, _ = json.Marshal(cfg)

	go log_file_date()

	go dial_server(dial_addr, mycfg, logHandler, shardConfig) // network.go

	wg.Add(2)
}

func start_full(dial_addr string) {

	fmt.Println("starting router...")
	start_router(dial_addr)
	time.Sleep(1 * time.Second)
	fmt.Println("starting log server...")
	start_log(config.Host + ":" + config.Port)
	time.Sleep(1 * time.Second)
	fmt.Println("starting shard...")
	start_shard(config.Host + ":" + config.Port)
	fmt.Println("starting client...")
	start_client(config.Host + ":" + config.Port)
}
