package native

import (
	"net"

	"github.com/dop251/goja"
)

func NetModuleLoader(vm *goja.Runtime, module *goja.Object) {
	export := module.Get("export").ToObject(vm)

	export.DefineDataProperty("BlockList", vm.ToValue(
		func(c goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
			b := &BlockList{
				vm:   vm,
				list: []SocketAddress{},
			}
			o := vm.NewDynamicObject(b)
			o.SetPrototype(c.This.Prototype())
			return o
		}), 1, 1, 1)
}

var _ goja.DynamicObject = (*BlockList)(nil)

type BlockList struct {
	vm   *goja.Runtime
	list []SocketAddress
}

func (b *BlockList) Get(key string) goja.Value {
	switch key {
	case "addAdress":
		return b.vm.ToValue(
			func(arg goja.FunctionCall, vm *goja.Runtime) goja.Value {
				if len(arg.Arguments) > 0 {

					switch v := arg.Arguments[0].Export().(type) {
					case string:
						ip := "ip4"
						if len(arg.Arguments) > 1 {
							if arg.Arguments[1].String() == "ipv6" {
								ip = "ip6"
							}
						}
						addr, _ := net.ResolveIPAddr(ip, v)
						return vm.NewDynamicObject(&SocketAddress{
							addr: addr.IP,
							zone: addr.Zone,
						})
					case *SocketAddress:
						b.list = append(b.list, *v)
					}
				}

				return goja.Undefined()
			})
	}
	return goja.Undefined()
}

var _ goja.DynamicObject

type SocketAddress struct {
	addr net.IP
	zone string
	port string
}

func (s *SocketAddress) Get(key string) goja.Value {}
