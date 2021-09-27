package config

import (
	"strconv"

	"rogchap.com/v8go"
)

func JsHandle(ctx *v8go.Context, uid, gid uint32, args ...string) (*v8go.Value, error) {
	switch args[0] {
	case "ML_enabled":
		return ctx.RunScript(strconv.FormatBool(Enable_ML), "cfg.js")
	case "enable_ML":
		Enable_ML = true

	case "disable_ML":
		Enable_ML = false

	case "app_port":

		return ctx.RunScript("'"+App_port+"'", "cfg.js")
	case "client_port":
		return ctx.RunScript("'"+Client_port+"'", "cfg.js")

	case "DebugMode":
		return ctx.RunScript(strconv.FormatBool(Debug), "cfg.js")

	case "autoShard":
		return ctx.RunScript(strconv.FormatBool(Auto_shard), "cfg.js")

	case "jstimeout":
		return ctx.RunScript(strconv.Itoa(Javascript_timeout), "cfg.js")

	case "chjstimeout":
		a, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}
		Javascript_timeout = a
	}
	return nil, nil
}
