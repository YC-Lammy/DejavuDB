package fmt

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

var _ require.ModuleLoader = FmtModuleLoader

func FmtModuleLoader(vm *goja.Runtime, module *goja.Object) {
	module.Set("Sprint", vm.ToValue(func(arg goja.FunctionCall) goja.Value {
		a := []interface{}{}
		for _, v := range arg.Arguments {
			a = append(a, v.Export())
		}
		return vm.ToValue(fmt.Sprint(a...))
	}))

	module.Set("Sprintf", vm.ToValue(func(arg goja.FunctionCall) goja.Value {
		a := []interface{}{}
		for _, v := range arg.Arguments[1:] {
			a = append(a, v.Export())
		}
		return vm.ToValue(fmt.Sprintf(arg.Arguments[0].String(), a...))
	}))

	module.Set("Sprintln", vm.ToValue(func(arg goja.FunctionCall) goja.Value {
		a := []interface{}{}
		for _, v := range arg.Arguments {
			a = append(a, v.Export())
		}
		return vm.ToValue(fmt.Sprintln(a...))
	}))

}
