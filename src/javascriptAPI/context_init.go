package javascriptAPI

import (
	"rogchap.com/v8go"

	_ "embed"
)

//go:embed api.js
var javascript_API_Script string

//go:embed adm.js
var javascript_ADM_Script string

func Javascript_context_init(vm *v8go.Isolate, errs chan error, delay_fn chan *func(), tmp_store map[string]interface{}, mode string, args ...[2]string) *v8go.Context {

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

	arg := map[string]string{"gid": gid, "uid": uid}

	glob, err := v8go.NewObjectTemplate(vm)
	if err != nil {
		errs <- err
		return nil
	}

	call_go_fn, err := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		return callbackfn(info, errs, delay_fn, arg, tmp_store)
	})
	if err != nil {
		errs <- err
		return nil
	}

	glob.Set("call_go_fn", call_go_fn) // register function

	ctx, err := v8go.NewContext(vm, glob)
	if err != nil {
		errs <- err
		return nil
	}

	ctx.RunScript(javascript_API_Script, "DB.js")

	if adm {
		ctx.RunScript(javascript_ADM_Script, "ADM.js")
	}

	switch mode {
	case "adm":
		ctx.RunScript(javascript_ADM_Script, "ADM.js")
	}

	return ctx
}
