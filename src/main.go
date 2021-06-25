package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"sync"
)

var password string = ""

var role string = ""

var mycfg = []byte{}

var MAC_Address string = get_first_mac_addr() // get.go

var hostport string = ""

var dial_ip string

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

	if password != "a empty password" {
		if len(password) != 16 && len(password) != 24 && len(password) != 32 {
			panic("password must be length of 16, 24 or 32")
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

}
