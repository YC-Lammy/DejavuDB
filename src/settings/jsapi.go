package settings

import "strconv"

func JsHandle(args ...string) (interface{}, error) {
	switch args[0] {
	case "ML_enabled":
		return Enable_ML, nil
	case "enable_ML":
		Enable_ML = true

	case "disable_ML":
		Enable_ML = false

	case "app_port":
		return App_port, nil
	case "client_port":
		return Client_port, nil

	case "DebugMode":
		return Debug, nil

	case "autoShard":
		return Auto_shard, nil

	case "jstimeout":
		return Javascript_timeout, nil

	case "chjstimeout":
		a, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		Javascript_timeout = a
	}
	return nil, nil
}
