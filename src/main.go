package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"src/lazy"
	"src/settings"

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

	fmt.Println("save to disk: ", settings.Save_disk)

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

	if settings.Password != "a empty password" {
		for {
			if len(settings.Password) != 16 && len(settings.Password) != 24 && len(settings.Password) != 32 {
				fmt.Println("password must be length of 16, 24 or 32")
				fmt.Scanln(&settings.Password)
			} else {
				break
			}
		}

	}

	fmt.Println("role: " + settings.Role + " listener ip: " + settings.Router_addr)
	switch settings.Role {

	case "router":
		start_router(settings.Router_addr)

	case "shard":
		start_shard(settings.Router_addr)

	case "client":
		start_client(settings.Router_addr)

	case "full":
		start_full(settings.Router_addr)

	case "log":
		start_log(settings.Router_addr)

	default:
		panic("Specified Role Invalid")

	}

	wg.Wait() // wait until all worker end

}

//////////////////////////////////////////////////////////////////////////////////////////

func start_router(dial_addr string) { // start as a router

	cfg := map[string]interface{}{"role": "router", "pass": settings.Password, "mac": MAC_Address, "port": settings.Host + ":" + settings.Port}
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
	cfg := map[string]interface{}{"role": "shard", "pass": settings.Password, "mac": MAC_Address, "port": settings.Host + ":" + settings.Port}
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
	cfg := map[string]interface{}{"role": "log", "pass": settings.Password, "port": settings.Host + ":" + settings.Port}
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
	start_log(settings.Host + ":" + settings.Port)
	time.Sleep(1 * time.Second)
	fmt.Println("starting shard...")
	start_shard(settings.Host + ":" + settings.Port)
	fmt.Println("starting client...")
	start_client(settings.Host + ":" + settings.Port)
}
