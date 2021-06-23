package router

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"log"
)

var password string

var role string

var MAC_Address string = get_first_mac_addr()

func main() {
	var router_addr string

	// flags declaration using flag package
	flag.StringVar(&role, "r", "router", "Specify role. Default is router, option: shard, client")
	flag.StringVar(&router_addr, "a", "", "Specify router ip. Default is stand alone router")
	flag.StringVar(&password, "p", "", "Specify pass. Default is empty")
	flag.Parse()
	gob.Register(map[string]interface{}{})
	switch role {
	case "router":
		start_router(router_addr)

	case "shard":
		if router_addr == "" {
			log.Fatalln("must specific an address")
			return
		}
		start_shard(router_addr)

	case "client":
		start_client(router_addr)

	default:
		log.Fatalln("Specified Role Invalid")

	}

}

func start_router(dial_addr string) {

	go start_listening() // router.go

	cfg := map[string]interface{}{"role": "router", "pass": password}
	mycfg, _ := json.Marshal(cfg)

	if dial_addr != "" {
		dial_server(dial_addr, string(mycfg), RouterHandler) // network.go
	}

}

func start_shard(dial_addr string) {

}

func start_client(dial_addr string) {

}
