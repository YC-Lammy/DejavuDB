package main

import (
	"errors"
	"fmt"

	"rogchap.com/v8go"
)

const javascript_API_Script = `
class dejavu_api_settings {
	constructor(){}

	static is_ML_enabled(){
		return dejavu_api_is_ML_enabled()
	}

	static enable_ML(){
		dejavu_api_enable_ML()
	}

	static disable_ML(){
		
	}
}

class dejavu_api_class {
	constructor(){

	}
	static version = "dejavuDB 0.2.1"
	static version_info = ""

	static settings = dejavu_api_settings
}
`

func javascript_context_init(ctx *v8go.Context, errs chan error) {

	vm, _ := ctx.Isolate()
	glob := ctx.Global()

	includefn, _ := v8go.NewFunctionTemplate(vm,
		func(info *v8go.FunctionCallbackInfo) *v8go.Value { // when the JS function is called this Go callback will execute

			switch v := fmt.Sprintf("%v", info.Args()); v {
			case `"dejavu"`, `dejavu`, `'dejavu'`, `"dejavu.js"`, `dejavu.js`, `'dejavu.js'`:
				ctx.RunScript(javascript_API_Script, "dejavuDB.js")
				val, _ := ctx.Global().Get("dejavu_api_class")
				return val

			default:
				errs <- errors.New("could not import" + v + " (cannot find package " + v + " in any of path)")

			}
			return nil // you can return a value back to the JS caller if required
		})

	glob.Set("include", includefn) // register function

	dejavu_api_is_ML_enabled, _ := v8go.NewFunctionTemplate(vm,

		func(info *v8go.FunctionCallbackInfo) *v8go.Value {

			v, _ := v8go.NewValue(vm, Settings.enable_ML)
			return v
		})

	glob.Set("dejavu_api_is_ML_enabled", dejavu_api_is_ML_enabled)

	dejavu_api_enable_ML, _ := v8go.NewFunctionTemplate(vm,

		func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			Settings.enable_ML = true
			return nil
		})

	glob.Set("dejavu_api_enable_ML", dejavu_api_enable_ML)
}
