package javascriptAPI

import (
	"rogchap.com/v8go"

	_ "embed"
)

//go:embed api.js
var javascript_API_Script string

//go:embed adm.js
var javascript_ADM_Script string

func Javascript_context_init(ctx *v8go.Context, errs chan error, delay_fn chan *func(), args ...[2]string) {

	adm := false

	var uid = "19890604"
	var gid = "19890604"

	for _, v := range args {
		switch v[0] {
		case "adm":
			adm = true
		case "gid":
			gid = v[1]
		case "uid":
			uid = v[1]
		}
	}

	arg_arr := [][2]string{
		[2]string{"gid", gid},
		[2]string{"uid", uid},
	}

	vm, _ := ctx.Isolate()
	glob := ctx.Global()

	call_go_fn, _ := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		return callbackfn(info, ctx, errs, delay_fn, arg_arr...)
	})

	glob.Set("call_go_fn", call_go_fn) // register function

	ctx.RunScript(javascript_API_Script, "DB.js")

	if adm {
		ctx.RunScript(javascript_ADM_Script, "ADM.js")
	}
}
