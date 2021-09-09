package javascriptAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"rogchap.com/v8go"

	"src/settings"
	tf "src/tensorflow"

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

var javascript_vm *v8go.Isolate

func init() {

}

func Javascript_context_init(ctx *v8go.Context, errs chan error, delay_fn chan *func()) {

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

			case "http_CanonicalHeaderKey":
				val, _ := v8go.NewValue(vm, http.CanonicalHeaderKey(Args[1].String()))
				return val

			case "http_DetectContentType":
				val, _ := v8go.NewValue(vm, http.DetectContentType([]byte(Args[1].String())))
				return val

			case "http_Get", "http_Head", "http_Post", "http_PostForm":
				type result struct {
					Status     string // e.g. "200 OK"
					StatusCode int    // e.g. 200
					Proto      string // e.g. "HTTP/1.0"
					ProtoMajor int    // e.g. 1
					ProtoMinor int    // e.g. 0

					Header map[string][]string
					Body   []byte

					ContentLength    int64
					TransferEncoding []string
					Uncompressed     bool

					Trailer map[string][]string
				}
				var resp *http.Response
				var err error
				switch Args[0].String() {
				case "http_Get":
					resp, err = http.Get(Args[1].String())
				case "http_Head":
					resp, err = http.Head(Args[1].String())
				case "http_Post", "http_PostForm":

					f, _ := os.CreateTemp("", "")
					defer f.Close()
					var ctype string

					switch Args[0].String() {
					case "PostForm":
						f.Write([]byte(Args[2].String()))
						ctype = "application/json; charset=UTF-8"
					case "Post":
						f.Write([]byte(Args[3].String()))
						ctype = Args[2].String()
					}

					resp, err = http.Post(Args[1].String(), ctype, f)

				}
				defer resp.Body.Close()

				if err != nil {
					errs <- err
					return nil
				}

				body, err := io.ReadAll(resp.Body)
				r := result{Status: resp.Status, StatusCode: resp.StatusCode, Proto: resp.Proto,
					ProtoMajor: resp.ProtoMajor, ProtoMinor: resp.ProtoMinor, Header: resp.Header,
					Body: body, ContentLength: resp.ContentLength, TransferEncoding: resp.TransferEncoding,
					Uncompressed: resp.Uncompressed, Trailer: resp.Trailer}

				barr, _ := json.Marshal(r)

				val, err := v8go.JSONParse(ctx, string(barr))
				if err != nil {
					errs <- err
					return nil
				}

				return val

			}
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

	//ctx.RunScript(javascript_API_Script, "dejavuDB.js")
}
