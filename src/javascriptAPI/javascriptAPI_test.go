package javascriptAPI

import (
	"testing"
)

func test_application_api(t *testing.T) {
	Javascript_run_isolate(`
DB.set('test.string','hello world',)
fmt = require('fmt')
fmt.Println(DB.Get('test.string'))
`)
}
