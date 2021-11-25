package javascriptAPI

import (
	"testing"
)

func test_application_api(t *testing.T) {
	JavascriptRun(`
DB.set('test.string','hello world',)
fmt = require('fmt')
fmt.Println(DB.Get('test.string'))
`, JsOptions{
		UserName: "",
	})
}
