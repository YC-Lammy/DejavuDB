package javascriptAPI

import "sync"

var javascript_API_lib = map[string]javascript_module{}
var javascript_API_lib_lock = sync.Mutex{}

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

func NewLib(name, version, version_info, auther, desc, js string) {

}
