package javascriptAPI

import (
	"dejavuDB/src/auth"
	"dejavuDB/src/datastore"

	"github.com/dop251/goja"
)

var _ goja.DynamicObject = (*StorageInterface)(nil)

type StorageInterface struct {
	Sets    map[string]interface{}
	Deletes []string

	user string
	vm   *goja.Runtime
}

func (s *StorageInterface) Get(key string) goja.Value {
	switch key {
	case "get":
		return s.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			if len(arg.Arguments) > 0 {
				return s.Get(arg.Arguments[0].String())
			}
			return goja.Undefined()
		})
	case "set":
		return s.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			if len(arg.Arguments) > 1 {
				s.Set(arg.Arguments[0].String(), arg.Arguments[1])
			}
			return goja.Undefined()
		})
	case "delete":
		return s.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			if len(arg.Arguments) > 0 {
				s.Delete(arg.Arguments[0].String())
			}
			return goja.Undefined()
		})
	case "copy":
		return s.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			if len(arg.Arguments) > 0 {
				return s.Get(arg.Arguments[0].String())
			}
			return goja.Undefined()
		})
	case "move":
		return s.vm.ToValue(func(arg goja.FunctionCall) goja.Value {
			if len(arg.Arguments) > 1 {
				v := s.Get(arg.Arguments[0].String())
				s.Set(arg.Arguments[1].String(), v)
			}
			return goja.Undefined()
		})
	default:
		if auth.HasPermission(s.user, key) {
			return s.vm.ToValue(datastore.Get(key))
		}
	}
	return goja.Undefined()
}

func (s *StorageInterface) Set(key string, val goja.Value) bool {
	if auth.HasPermission(s.user, key) {
		s.Sets[key] = val.Export()
		return true
	}
	return false
}

func (s *StorageInterface) Delete(key string) bool {
	if auth.HasPermission(s.user, key) {
		s.Deletes = append(s.Deletes, key)
		return true
	}
	delete(s.Sets, key)
	return false
}

func (s *StorageInterface) Has(key string) bool {
	return auth.HasPermission(s.user, key)
}

func (s *StorageInterface) Keys() []string {
	return []string{}
}

func (s *StorageInterface) Commit()
