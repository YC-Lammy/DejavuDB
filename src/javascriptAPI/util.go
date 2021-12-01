package javascriptAPI

import (
	"strings"

	"github.com/dop251/goja"
)

var jsInstanceofWrapper = goja.MustCompile("", `
function(a, b){
	return a instanceof b
}
`, false)

func InstanceOf(vm *goja.Runtime, target goja.Value, Type string) bool {
	s := strings.Split(Type, ".")
	dst := vm.GlobalObject().Get(s[0])
	for _, v := range s[1:] {
		dst = dst.ToObject(vm).Get(v)
	}
	fn, _ := vm.RunProgram(jsInstanceofWrapper)
	c, _ := goja.AssertFunction(fn)
	b, _ := c(goja.Null(), target, dst)
	return b.ToBoolean()
}
