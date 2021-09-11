package plugins

import (
	"../javascriptAPI"
	"rogchap.com/v8go"
)

var plugin_vm *v8go.Isolate

var plugin_register = map[string]plugin{}

// plugins written in javascript

type plugin struct {
	plugin_name            string
	plugin_version         string
	plugin_status          bool
	plugin_type            string // "type", "function", "service"
	plugin_type_version    string
	plugin_library         string
	plugin_library_version string
	plugin_author          string
	plugin_description     string
	plugin_license         string
	load_option            string
	plugin_maturity        string
	plugin_auth_version    string

	script          string
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

func (plug plugin) NewContext() (string, error) {
	{
		errs := make(chan error, 1)
		delay_fn := make(chan *func(), 1)

		ctx, err := v8go.NewContext()

		if err != nil {
			errs <- err
			return "", err
		}
		defer ctx.Close()
		vm, _ := ctx.Isolate()

		defer vm.TerminateExecution()

		javascriptAPI.Javascript_context_init(ctx, errs, delay_fn) // initiallize context api and functions

		val, err := ctx.RunScript(plug.script, "main.js") // exec a long running script

		return val.String(), nil

	}
}

func register_plugin() {

}
