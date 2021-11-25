package javascriptAPI

import (
	"dejavuDB/src/config"

	"github.com/dop251/goja"
	"rogchap.com/v8go"
)

type Insruction struct {
	Insruction string
}

var _ goja.DynamicObject = (*DatabaseObject)(nil)

type DatabaseObject struct {
	variable map[string]goja.Value
	vm       *goja.Runtime

	storage *StorageInterface
}

func NewVmDatabase(vm *goja.Runtime) *DatabaseObject {
	db := &DatabaseObject{
		variable: make(map[string]goja.Value),
		vm:       vm,
	}
	vm.GlobalObject().DefineDataProperty("database", vm.NewDynamicObject(db), 1, 1, 1)
	vm.Set("database", vm.NewDynamicObject(db))
	return db
}

func (d *DatabaseObject) Set(key string, val goja.Value) bool {
	d.variable[key] = val
	return true
}

func (d *DatabaseObject) Delete(key string) bool {
	delete(d.variable, key)
	return true
}

func (d *DatabaseObject) Has(key string) bool {
	for _, v := range d.Keys() {
		if v == key {
			return true
		}
	}
	if _, ok := d.variable[key]; ok {
		return true
	}
	return false
}

func (d *DatabaseObject) Keys() []string {
	return []string{
		"Storage",
	}
}

func (d *DatabaseObject) Get(key string) goja.Value {

	switch key {
	case "Storage":
		switch config.Role {
		case "router":

		case "standalone":
			return d.vm.ToValue(d.vm.NewDynamicObject(d.storage))
		}
	}
	return nil
}

func checkerr(err error, val *v8go.Value, errs chan error) *v8go.Value {
	if err != nil {
		errs <- err
		return nil
	}
	return val
}
