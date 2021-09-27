package javascriptAPI

import (
	"time"

	"rogchap.com/v8go"

	"src/config"
)

func Javascript_run_isolate(vm *v8go.Isolate, script string, mode string, args ...interface{}) (string, error) {
	vals := make(chan *v8go.Value, 1)
	errs := make(chan error, 1)
	delay_fn := make(chan *func(), 1)
	delay_fns := []*func(){}

	var ctx *v8go.Context

	tmp_store := map[string]interface{}{}

	go func() {

		ctx = Javascript_context_init(vm, errs, delay_fn, tmp_store, mode) // initiallize context api and functions

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

		val, err := ctx.RunScript("return_print_buffer", "return.js")
		if err != nil {
			return "", err
		}
		s := val.String()
		ctx.Close()
		return s, nil
	case err := <-errs:
		// javascript error
		for _, fn := range delay_fns {
			(*fn)()
		}
		ctx.Close()
		return "", err
	case <-time.After(time.Duration(config.Javascript_timeout) * time.Millisecond): // get the Isolate from the context
		for _, fn := range delay_fns {
			(*fn)()
		}
		ctx.Close()
		err := <-errs // will get a termination error back from the running script
		return "", err
	}
}
