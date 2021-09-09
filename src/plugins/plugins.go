package plugins

import "C"
import "rogchap.com/v8go"

var plugin_vm *v8go.Isolate

var plugin_register = map[string]plugin{}
// plugins written in javascript

type plugin struct {
	plugin_name            string
	plugin_version         string
	plugin_status          bool
	plugin_type            string // "type", "function"
	plugin_type_version    string
	plugin_library         string
	plugin_library_version string
	plugin_author          string
	plugin_description     string
	plugin_license         string
	load_option            string
	plugin_maturity        string
	plugin_auth_version    string

	command         string
	execute         func(...interface{})
	result_key_word string
	error_terminate bool
}

func init() {
	vm, err := v8go.NewIsolate()
	if err != nil {
		panic(err)
	}
	plugin_vm = vm
}

func NewContext() (string, error){
	{
		vals := make(chan *v8go.Value, 1)
		errs := make(chan error, 1)
		delay_fn := make(chan *func(), 1)
		delay_fns := []*func(){}
		ctx, err := v8go.NewContext(plugin_vm)
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
	
			javascript_context_init(ctx, errs, delay_fn) // initiallize context api and functions
	
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
}

func register_plugin() {

}
