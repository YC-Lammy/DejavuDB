package main

func shardConfig(config map[string]interface{}) error {

	if v, ok := config["router_ipv4"]; ok {

		var list []string
		for _, v := range v.([]interface{}) { // convert []interface{} to []string
			list = append(list, v.(string))
		}
		current_router_ipv4 = removeDuplicateStrings(list)

		for _, ip := range current_router_ipv4 {
			if ip != dial_ip {
				go dial_server(ip, mycfg, ShardHandler, secondConfig)
				wg.Add(1)
			}

		}
	}
	return nil
}

func routerConfig(config map[string]interface{}) error {
	if v, ok := config["router_ipv4"]; ok {

		var list []string
		for _, v := range v.([]interface{}) { // convert []interface{} to []string
			list = append(list, v.(string))
		}

		current_router_ipv4 = removeDuplicateStrings(list)
		for _, ip := range current_router_ipv4 {
			if ip != dial_ip && ip != hostport {
				go dial_server(ip, mycfg, ShardHandler, secondConfig)
				wg.Add(1)
			}

		}
	}
	return nil
}

func secondConfig(config map[string]interface{}) error { // secondary configaration, connect to all routers
	if v, ok := config["router_ipv4"]; ok {
		var list []string
		for _, v := range v.([]interface{}) { // convert []interface{} to []string
			list = append(list, v.(string))
		}

		more_ip := difference(current_router_ipv4, list)

		more_ip = removeDuplicateStrings(more_ip)

		for _, ip := range more_ip {

			if ip != dial_ip {
				go dial_server(ip, mycfg, ShardHandler, secondConfig)
				wg.Add(1)
			}

			current_router_ipv4 = append(current_router_ipv4, ip)
		}

	}
	return nil
}
