package javascriptAPI

import (
	"testing"

	"rogchap.com/v8go"
)

func test_application_api(t *testing.T) {
	iso, _ := v8go.NewIsolate()
	Javascript_run_isolate(iso, `
DB.set('test.string','hello world',)
fmt = require('fmt')
fmt.Println(DB.Get('test.string'))
`, "")
}
