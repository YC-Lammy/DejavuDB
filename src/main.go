package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var password string = ""

var role string = ""

var mycfg = []byte{}

var MAC_Address string = get_first_mac_addr() // get.go

var hostport string = ""

var dial_ip string = ""

var save_to_disk bool = false

var securite_connection int = 0

var ML = false

var sql_file string = ""

var DEBUG bool = false

var wg sync.WaitGroup // working group

func main() {
	var router_addr string

	// flags declaration using flag package
	flag.StringVar(&role, "r", "router", "Specify role. Default is router, option: shard, client")
	flag.StringVar(&router_addr, "ip", "", "Specify router ip. Default is stand alone router")
	flag.StringVar(&password, "p", "a empty password", "Specify password. Default is empty")
	flag.StringVar(&hostport, "host", "localhost:8080", "specify hosting port")
	flag.StringVar(&sql_file, "sqlfile", ":memory:", "specify sql file path")
	flag.BoolVar(&save_to_disk, "disk", false, "save copy to disk")
	flag.BoolVar(&DEBUG, "debug", false, "developer debug option")
	flag.BoolVar(&ML, "Machine Learning", false, "Enable built in machine learning service")
	flag.IntVar(&securite_connection, "sc", 0, "specify to use securite connection and the bit width")
	flag.Parse()
	gob.Register(map[string]interface{}{})

	//fmt.Println("enter your password:")

	//fmt.Scanln(&password)

	fmt.Println("save to disk: ", save_to_disk)

	os.Chdir(home_dir)

	os.Mkdir("dejavuDB", os.ModePerm)

	os.Chdir("dejavuDB")

	os.Mkdir("log", os.ModePerm)
	os.Mkdir("database", os.ModePerm)
	os.Mkdir("ML", os.ModePerm)
	os.Chdir("database")
	if _, err := os.Stat("dejavu.db"); os.IsNotExist(err) {
		f, _ := os.Create("dejavu.db")
		f.Close()

	}
	sql_file = filepath.Join(home_dir, "dejavuDB", "dejavu.db")
	os.Chdir(path.Join(home_dir, "dejavuDB"))

	setupLog()

	if password != "a empty password" {
		for {
			if len(password) != 16 && len(password) != 24 && len(password) != 32 {
				fmt.Println("password must be length of 16, 24 or 32")
				fmt.Scanln(&password)
			} else {
				break
			}
		}

	}

	dial_ip = router_addr

	fmt.Println("role: " + role + " listener ip: " + router_addr)
	switch role {

	case "router":
		start_router(router_addr)

	case "shard":
		start_shard(router_addr)

	case "client":
		start_client(router_addr)

	case "full":
		start_full(router_addr)

	case "log":
		start_log(router_addr)

	default:
		panic("Specified Role Invalid")

	}

	wg.Wait() // wait until all worker end

}

//////////////////////////////////////////////////////////////////////////////////////////

func start_router(dial_addr string) { // start as a router

	cfg := map[string]interface{}{"role": "router", "pass": password, "mac": MAC_Address, "port": hostport}
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
	cfg := map[string]interface{}{"role": "shard", "pass": password, "mac": MAC_Address, "port": hostport}
	mycfg, _ = json.Marshal(cfg)

	go dial_server(dial_addr, mycfg, ShardHandler, shardConfig) // network.go

	wg.Add(1)
}

////////////////////////////////////////////////////////////////////////////////////////

func start_client(dial_addr string) { // start as a client
	cfg := map[string]interface{}{"role": "client"}
	mycfg, _ = json.Marshal(cfg)
	go Client_dial(dial_addr, mycfg)

	go SQL_init()

	wg.Add(1)
}

func start_log(dial_addr string) {
	if dial_addr == "" {
		panic("must specific an address")
		return
	}
	cfg := map[string]interface{}{"role": "log", "pass": password, "port": hostport}
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
	start_log(hostport)
	time.Sleep(1 * time.Second)
	fmt.Println("starting shard...")
	start_shard(hostport)
	fmt.Println("starting client...")
	start_client(hostport)
}
