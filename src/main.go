package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

var mycfg = []byte{}

var MAC_Address string = get_first_mac_addr() // get.go

var sql_file string = ""

var wg sync.WaitGroup // working group

func main() {

	// flags declaration using flag package
	init_settings() // handles all flags

	for i, v := range os.Args {
		switch v {
		case "help":
			fmt.Println(manual_desc)
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

	fmt.Println("save to disk: ", Settings.save_disk)

	os.Chdir(home_dir)

	os.Mkdir("dejavuDB", os.ModePerm)

	os.Chdir("dejavuDB")

	os.Mkdir("log", os.ModePerm)
	os.Mkdir("database", os.ModePerm)
	os.Mkdir("ML", os.ModePerm)
	os.Mkdir("addon", os.ModePerm)

	init_users()

	os.Chdir(path.Join(home_dir, "dejavuDB"))

	os.Chdir("database")

	os.Mkdir("tables", os.ModePerm)

	//sql_file = filepath.Join(home_dir, "dejavuDB", "dejavu.db")
	os.Chdir(path.Join(home_dir, "dejavuDB"))

	setupLog()

	if Settings.password != "a empty password" {
		for {
			if len(Settings.password) != 16 && len(Settings.password) != 24 && len(Settings.password) != 32 {
				fmt.Println("password must be length of 16, 24 or 32")
				fmt.Scanln(&Settings.password)
			} else {
				break
			}
		}

	}

	fmt.Println("role: " + Settings.role + " listener ip: " + Settings.router_addr)
	switch Settings.role {

	case "router":
		start_router(Settings.router_addr)

	case "shard":
		start_shard(Settings.router_addr)

	case "client":
		start_client(Settings.router_addr)

	case "full":
		start_full(Settings.router_addr)

	case "log":
		start_log(Settings.router_addr)

	default:
		panic("Specified Role Invalid")

	}

	wg.Wait() // wait until all worker end

}

//////////////////////////////////////////////////////////////////////////////////////////

func start_router(dial_addr string) { // start as a router

	cfg := map[string]interface{}{"role": "router", "pass": Settings.password, "mac": MAC_Address, "port": Settings.host + ":" + Settings.port}
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
	cfg := map[string]interface{}{"role": "shard", "pass": Settings.password, "mac": MAC_Address, "port": Settings.host + ":" + Settings.port}
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
	cfg := map[string]interface{}{"role": "log", "pass": Settings.password, "port": Settings.host + ":" + Settings.port}
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
	start_log(Settings.host + ":" + Settings.port)
	time.Sleep(1 * time.Second)
	fmt.Println("starting shard...")
	start_shard(Settings.host + ":" + Settings.port)
	fmt.Println("starting client...")
	start_client(Settings.host + ":" + Settings.port)
}
