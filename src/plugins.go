package main

import "C"
import "rogchap.com/v8go"

// plugins written in python or javascript

type plugin struct {
	plugin_name            string
	plugin_version         string
	plugin_status          bool
	plugin_type            string
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
	language        string // "javascript" / "python"
	execute         func(...interface{})
	result_key_word string
}

var jsvm *v8go.Isolate

func init_plugins() error {
	vm, err := v8go.NewIsolate()
	if err != nil {
		return err
	}
	jsvm = vm
	return nil
}

func register_plugin() {

}
