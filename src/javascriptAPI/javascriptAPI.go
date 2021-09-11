package javascriptAPI

import (
	"github.com/goccy/go-json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"rogchap.com/v8go"

	"../settings"

	tf "../tensorflow"

	_ "embed"
)

type javascript_module struct {
	name         string
	version      string
	version_info string
	auther       string
	describtion  string
	model_path   string

	is_in_ram bool
	script    string
	enabled   bool
}

//go:embed api.js
var javascript_API_Script string

var javascript_API_lib = map[string]*javascript_module{}
var javascript_API_lib_lock = sync.Mutex{}

func Javascript_context_init(ctx *v8go.Context, errs chan error, delay_fn chan *func()) {

	var terminate_fns = []func(){}

	defer func(){
		for _, v := range terminate_fns{
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
				javascript_API_lib_lock.Lock()
				module, ok := javascript_API_lib[v]
				javascript_API_lib_lock.Unlock()

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

			case "Set":

			case "dejavu_api_is_ML_enabled":
				v, _ := v8go.NewValue(vm, settings.Enable_ML)
				return v

			case "dejavu_api_enable_ML":
				settings.Enable_ML = true
				return nil

			case "dejavu_api_disable_ML":
				settings.Enable_ML = false
				return nil

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

			
			return nil // you can return a value back to the JS caller if required
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
