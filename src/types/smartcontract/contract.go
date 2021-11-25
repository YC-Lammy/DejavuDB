package contract

import (
	"strings"
	"time"

	"rogchap.com/v8go"
)

type Contract struct {
	LocalStore map[string]*v8go.Value
	Script     string
	Name       string
}

func (contr *Contract) Run(function string, args ...*v8go.Value) (string, error) {
	vals := make(chan *v8go.Value, 1)
	errs := make(chan error, 1)
	ctx, err := v8go.NewContext()

	defer ctx.Close()

	if err != nil {
		errs <- err
		return "", err
	}

	vm, _ := ctx.Isolate()

	go func() {

		_, err := ctx.RunScript(contr.Script, "contract.js") // exec a long running script
		if err != nil {
			errs <- err
			return "", nil
		}

		glob := ctx.Global()

		for key, val := range contr.LocalStore {
			err := glob.Set(key, val)
			if err != nil {
				errs <- err
				return "", nil
			}
		}

		v, err := glob.Get(function)
		if err != nil {
			errs <- err
			return "", nil
		}

		fn, err := v.AsFunction()
		if err != nil {
			errs <- err
			return "", nil
		}

		val, err := fn.Call(args...)
		if err != nil {
			errs <- err
			return "", nil
		}
		vals <- val
	}()

	select {
	case val := <-vals:
		// sucess
		keys, err := ctx.RunScript(`
		const keys = [];
		for (const [key, value] of Object.entries(object1)) {
			keys.push(key);
		  };
		String(keys);`, "storage.js")
		if err != nil {
			return "", err
		}

		glob := ctx.Global()
		tmp := map[string]*v8go.Value{}
		for _, k := range strings.Split(keys, ",") {
			v, err := glob.Get(k)
			if err != nil {
				return "", err
			}
			tmp[k] = v
		}

		contr.LocalStore = tmp
		return val.String(), nil

	case err := <-errs:
		// javascript error
		return "", err
	case <-time.After(time.Duration(1000) * time.Millisecond): // get the Isolate from the context
		vm.TerminateExecution() // terminate the execution
		err := <-errs           // will get a termination error back from the running script
		return "", err
	}
}
