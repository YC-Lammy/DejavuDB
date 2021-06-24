package main

func shardConfig(config map[string]interface{}) error {
	if v, ok := config["router_ipv4"]; ok {
		current_router_ipv4 = v.([]string)
		for _, ip := range v.([]string) {
			go dial_server(ip, mycfg, ShardHandler, secondConfig)
			wg.Add(1)
		}
	}
	return nil
}

func routerConfig(config map[string]interface{}) error {
	if v, ok := config["router_ipv4"]; ok {
		current_router_ipv4 = v.([]string)
		for _, ip := range v.([]string) {
			go dial_server(ip, mycfg, ShardHandler, secondConfig)
			wg.Add(1)
		}
	}
	return nil
}

func secondConfig(config map[string]interface{}) error {
	if v, ok := config["router_ipv4"]; ok {
		more_ip := difference(current_router_ipv4, v.([]string))

		for _, ip := range more_ip {

			go dial_server(ip, mycfg, ShardHandler, secondConfig)
			wg.Add(1)

			current_router_ipv4 = append(current_router_ipv4, ip)
		}

	}
	return nil
}
