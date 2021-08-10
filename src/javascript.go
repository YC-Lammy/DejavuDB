package main

import (
	"time"

	"rogchap.com/v8go"
)

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

		javascript_context_init(ctx, errs) // initiallize context api and functions

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
