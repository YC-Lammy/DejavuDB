package javascriptAPI

import (
	"time"

	"src/settings"

	"rogchap.com/v8go"
)

func Javascript_run_isolate(script string, args ...interface{}) (string, error) {
	vals := make(chan *v8go.Value, 1)
	errs := make(chan error, 1)
	delay_fn := make(chan *func(), 1)
	delay_fns := []*func(){}
	ctx, err := v8go.NewContext()
	if err != nil {
		errs <- err
		return "", err
	}

	vm, _ := ctx.Isolate()

	go func() {
		fn := <-delay_fn
		delay_fns = append(delay_fns, fn)
	}()

	go func() {

		Javascript_context_init(ctx, errs, delay_fn) // initiallize context api and functions

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
		for _, fn := range delay_fns {
			a := *fn
			a()
		}
		return val.String(), nil
	case err := <-errs:
		// javascript error
		return "", err
	case <-time.After(time.Duration(settings.Javascript_timeout) * time.Millisecond): // get the Isolate from the context
		vm.TerminateExecution() // terminate the execution
		err := <-errs           // will get a termination error back from the running script
		return "", err
	}
}
