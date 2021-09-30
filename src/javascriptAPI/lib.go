package javascriptAPI

import (
	_ "embed"
	"sync"
)

var javascript_API_lib = map[string]javascript_module{}
var javascript_API_lib_lock = sync.RWMutex{}

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

//go:embed lib/http/http.js
var _http_js_script_ string

//go:embed lib/fmt/fmt.js
var _fmt_js_script_ string

func init() {
	javascript_API_lib["http"] = javascript_module{
		name:         "http",
		version:      "",
		version_info: "",
		auther:       "YC",
		describtion:  "javascript wrapper for go http",
		model_path:   "builtin",

		is_in_ram: true,
		script:    _http_js_script_,
		enabled:   false, // default disabled for security
	}

	javascript_API_lib["fmt"] = javascript_module{
		name:         "fmt",
		version:      "",
		version_info: "",
		auther:       "YC",
		describtion:  "javascript wrapper fo go fmt",
		model_path:   "builtin",

		is_in_ram: true,
		script:    _fmt_js_script_,
		enabled:   true,
	}
}
