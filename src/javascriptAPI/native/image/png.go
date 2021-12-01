package image

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/dop251/goja"
)

func ImagePNGModuleLoader(vm *goja.Runtime, module *goja.Object) {
	e := module.Get("export")
	export := e.ToObject(vm)

	export.Set("decode", vm.ToValue(func(arg goja.FunctionCall) goja.Value {
		b := bytes.NewBuffer([]byte{})
		im, f, err := image.Decode(b)
		if err != nil {
			vm.RunString("throw '" + err.Error() + "'")
		}
		return vm.NewDynamicObject(&Image{
			vm:           vm,
			image:        im,
			encodeformat: f,
		})
	}))
}
