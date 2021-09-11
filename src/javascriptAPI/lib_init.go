package javascriptAPI

import (
	_ "embed"
)

//go:embed lib/http/http.js
var _http_js_script_ string

func init() {
	javascript_API_lib["http"] = &javascript_module{
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
}
