package image

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func Test1(t *testing.T) {
	re := require.NewRegistry()
	re.RegisterNativeModule("image", ImageModuleLoader)
	vm := goja.New()
	re.Enable(vm)

	_, err := vm.RunString(`
	image = require('image')
	im = image.noiseImage(0,0)
	Object.keys(im)
	im.
	`)
	if err != nil {
		panic(err)
	}
}
