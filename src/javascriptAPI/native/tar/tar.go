package tar

import (
	"archive/tar"
	"bytes"
	"io"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

var _ require.ModuleLoader = TarModuleLoader

func TarModuleLoader(vm *goja.Runtime, module *goja.Object) {
	e := module.Get("exports")
	export := e.ToObject(vm)

	export.Set("Reader", func(c goja.ConstructorCall) *goja.Object {
		obj := vm.NewDynamicObject(&reader{
			vm: vm,
		})
		obj.SetPrototype(c.This.Prototype())
		return obj
	})

	export.Set("Writer", func(c goja.ConstructorCall) *goja.Object {
		obj := vm.NewDynamicObject(&writer{
			vm: vm,
		})
		obj.SetPrototype(c.This.Prototype())
		return obj
	})

}

type reader struct {
	vm *goja.Runtime
}

func (r reader) Get(key string) goja.Value {
	switch key {
	case "read":
		return r.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			if len(arg.Arguments) < 0 {
				return goja.Undefined()
			}
			var array goja.ArrayBuffer
			switch ar := arg.Arguments[0].Export().(type) {
			case goja.ArrayBuffer:
				array = ar
			default:
				if ar, ok := arg.Arguments[0].ToObject(r.vm).Get("buffer").Export().(goja.ArrayBuffer); ok {
					array = ar
				}
			}

			reader := tar.NewReader(bytes.NewReader(array.Bytes()))
			return r.vm.NewDynamicObject(result{
				Reader: reader,
			})
		})
	}
	return goja.Undefined()
}

func (r reader) Set(key string, val goja.Value) bool {
	return false
}

func (r reader) Delete(key string) bool {
	return false
}

func (r reader) Has(key string) bool {
	for _, k := range r.Keys() {
		if k == key {
			return true
		}
	}
	return false
}

func (r reader) Keys() []string {
	return []string{"read"}
}

type result struct {
	vm *goja.Runtime
	*tar.Reader
}

func (r result) Get(key string) goja.Value {
	switch key {
	case "next":
		next, err := r.Next()
		if err != nil {
			return goja.Undefined()
		}
		b, _ := io.ReadAll(r)
		return r.vm.NewDynamicObject(header{
			vm:     r.vm,
			data:   b,
			Header: next,
		})
	}
	return goja.Undefined()
}

func (r result) Set(key string, val goja.Value) bool {
	return false
}

func (r result) Delete(key string) bool {
	return false
}

func (r result) Has(key string) bool {
	for _, k := range r.Keys() {
		if k == key {
			return true
		}
	}
	return false
}

func (r result) Keys() []string {
	return []string{"next"}
}

type header struct {
	vm   *goja.Runtime
	data []byte
	*tar.Header
}

func (r header) Get(key string) goja.Value {
	switch key {
	case "name":
		return r.vm.ToValue(r.Name)
	case "format":
		return r.vm.ToValue(r.Format.String())
	case "devmajor":
		return r.vm.ToValue(r.Devmajor)
	case "devminor":
		return r.vm.ToValue(r.Devminor)
	case "linkname":
		return r.vm.ToValue(r.Linkname)
	case "size":
		return r.vm.ToValue(r.Size)
	case "typeflag":
		var t string
		switch r.Typeflag {
		case tar.TypeBlock:
			t = "block"
		case tar.TypeChar:
			t = "char"
		case tar.TypeCont:
			t = "cont"
		case tar.TypeDir:
			t = "dir"
		case tar.TypeFifo:
			t = "fifo"
		case tar.TypeGNULongLink:
			t = "gnuLongLink"
		case tar.TypeGNULongName:
			t = "gnuLongName"
		case tar.TypeGNUSparse:
			t = "gnuSparse"
		case tar.TypeLink:
			t = "link"
		case tar.TypeReg:
			t = "reg"
		case tar.TypeRegA:
			t = "rega"
		case tar.TypeSymlink:
			t = "symlink"
		case tar.TypeXGlobalHeader:
			t = "XGlobalHeader"
		case tar.TypeXHeader:
			t = "XHeader"
		}
		return r.vm.ToValue(t)
	case "data":
		return r.vm.ToValue(r.vm.NewArrayBuffer(r.data))
	}
	return goja.Undefined()
}

func (r header) Set(key string, val goja.Value) bool {
	return false
}

func (r header) Delete(key string) bool {
	return false
}

func (r header) Has(key string) bool {
	for _, k := range r.Keys() {
		if k == key {
			return true
		}
	}
	return false
}

func (r header) Keys() []string {
	return []string{}
}

type writer struct {
	vm *goja.Runtime
}

func (r writer) Get(key string) goja.Value {
	switch key {
	case "write":
	}
	return goja.Undefined()
}

func (r writer) Set(key string, val goja.Value) bool {
	return false
}

func (r writer) Delete(key string) bool {
	return false
}

func (r writer) Has(key string) bool {
	for _, k := range r.Keys() {
		if k == key {
			return true
		}
	}
	return false
}

func (r writer) Keys() []string {
	return []string{}
}
