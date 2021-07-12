package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

var password string = ""

var role string = ""

var mycfg = []byte{}

var MAC_Address string = get_first_mac_addr() // get.go

var hostport string = ""

var dial_ip string = ""

var wg sync.WaitGroup // working group

func main() {
	var router_addr string

	// flags declaration using flag package
	flag.StringVar(&role, "r", "router", "Specify role. Default is router, option: shard, client")
	flag.StringVar(&router_addr, "ip", "", "Specify router ip. Default is stand alone router")
	flag.StringVar(&password, "p", "a empty password", "Specify password. Default is empty")
	flag.StringVar(&hostport, "port", "localhost:8080", "specify hosting port")
	flag.Parse()
	gob.Register(map[string]interface{}{})

	//fmt.Println("enter your password:")

	//fmt.Scanln(&password)

	os.Chdir(home_dir)

	os.Mkdir("dejavuDB", os.ModePerm)

	os.Chdir("dejavuDB")

	os.Mkdir("log", os.ModePerm)
	os.Mkdir("database", os.ModePerm)

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
	c := color.New(color.FgHiRed).Add(color.Bold)
	c.Println("\nListening at " + hostport + "\n")
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

	wg.Add(1)
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
