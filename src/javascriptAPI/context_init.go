package javascriptAPI

import (
	"errors"
	"fmt"
	"io/ioutil"
	"src/datastore"
	"src/types"

	"rogchap.com/v8go"

	"../settings"

	tf "../tensorflow"

	_ "embed"
)

//go:embed api.js
var javascript_API_Script string

func Javascript_context_init(ctx *v8go.Context, errs chan error, delay_fn chan *func()) {

	var terminate_fns = []func(){}

	defer func() {
		for _, v := range terminate_fns {
			v()
		}
	}()

	vm, _ := ctx.Isolate()
	glob := ctx.Global()

	call_go_fn, _ := v8go.NewFunctionTemplate(vm,
		func(info *v8go.FunctionCallbackInfo) *v8go.Value { // when the JS function is called this Go callback will execute

			Args := info.Args()
			switch Args[0].String() {
			case "require":
				v := Args[1].String()
				module, ok := javascript_API_lib[v]

				if ok {
					var body string
					if module.enabled { // some module such as networking is disabled by default
						if !module.is_in_ram {
							b, err := ioutil.ReadFile(module.model_path)
							if err != nil {
								errs <- err
								return nil
							}
							body = string(b)
						} else {
							body = module.script
						}

						val, _ := v8go.NewValue(vm, string(body))

						return val
					}

				} else { // module not exist
					errs <- errors.New("could not import" + v + " (cannot find package " + v + " in any of path)")
				}

			case "Get":
				dtype, p := datastore.Get(Args[1].String())
				if p == nil {
					return nil
				}
				switch dtype {
				case types.String:
					val, _ := v8go.NewValue(vm, *(*string)(p))
					return val
				case types.Int, types.Int64:
					val, _ := v8go.NewValue(vm, *(*int64)(p))
					return val
				case types.Int32:
					val, _ := v8go.NewValue(vm, *(*int32)(p))
					return val
				}

			case "Set":

			case "settings":
				a := []string{}
				for _, v := range Args[1:] {
					a = append(a, v.String())
				}
				settings.JsHandle(a...)

			case "tensorflow":

			case "TF_MODEL_EXIST":
				name := fmt.Sprintf("%v", Args[1])
				_, err := tf.Get_model_by_name(name)
				a := true
				if err != nil {
					a = false
				}
				val, err := v8go.NewValue(vm, a)
				if err != nil {
					errs <- err
				}
				return val

			case "http":

				// you can return a value back to the JS caller if required
			}
			return nil
		})

	glob.Set("call_go_fn", call_go_fn) // register function

	ctx.RunScript(`
	function require(name){
		script = call_go_fn("require", name);
		var exports = {};
		eval(script)
		return exports
	}
	`, "require.js")

	ctx.RunScript(javascript_API_Script, "DB.js")
}
