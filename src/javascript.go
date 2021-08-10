package main

import (
	"errors"
	"fmt"
	"time"

	"rogchap.com/v8go"
)

func javascript_context_init(ctx *v8go.Context, errs chan error) {

	vm, _ := ctx.Isolate()

	includefn, _ := v8go.NewFunctionTemplate(vm,
		func(info *v8go.FunctionCallbackInfo) *v8go.Value { // when the JS function is called this Go callback will execute

			switch v := fmt.Sprintf("%v", info.Args()); v {
			case `"dejavu"`, `dejavu`, `'dejavu'`:
				ctx.RunScript(javascript_API_Script, "dejavuDB.js")
				val, _ := ctx.Global().Get("dejavu_api_class")
				return val

			default:
				errs <- errors.New("could not import" + v + " (cannot find package " + v + " in any of path)")

			}
			return nil // you can return a value back to the JS caller if required
		})

	ctx.Global().Set("include", includefn) // register function
}

func javascript_run_isolate(script string, args ...interface{}) (string, error) {
	vals := make(chan *v8go.Value, 1)
	errs := make(chan error, 1)
	ctx, err := v8go.NewContext()
	if err != nil {
		errs <- err
		return "", err
	}

	vm, _ := ctx.Isolate()

	go func() {

		javascript_context_init(ctx, errs)

		val, err := ctx.RunScript(script, "main.js") // exec a long running script
		if err != nil {
			errs <- err
			return
		}
		vals <- val
	}()

	select {
	case val := <-vals:
		// sucess
		return val.String(), nil
	case err := <-errs:
		// javascript error
		return "", err
	case <-time.After(time.Duration(Settings.javascript_timeout) * time.Millisecond): // get the Isolate from the context
		vm.TerminateExecution() // terminate the execution
		err := <-errs           // will get a termination error back from the running script
		return "", err
	}
}
