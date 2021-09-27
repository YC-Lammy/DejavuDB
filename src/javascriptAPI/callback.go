package javascriptAPI

import (
	"errors"
	"io/ioutil"
	"reflect"

	"strconv"

	"src/config"
	"src/datastore"
	"src/types"
	"src/types/int128"

	"rogchap.com/v8go"
)

func callbackfn(info *v8go.FunctionCallbackInfo, errs chan error, delayfn chan *func(), args map[string]string, tmp_store map[string]interface{}) *v8go.Value { // when the JS function is called this Go callback will execute
	defer func() {
		if err := recover(); err != nil {
			errs <- err.(error)
		}
	}()
	ctx := info.Context()

	var uid uint32 = 19890604
	var gid uint32 = 19890604

	for k, v := range args {
		switch k {
		case "gid":
			a, _ := strconv.ParseUint(v, 10, 32)
			gid = uint32(a)
		case "uid":
			a, _ := strconv.ParseUint(v, 10, 32)
			uid = uint32(a)
		}
	}
	vm, err := ctx.Isolate()
	if err != nil {
		errs <- err
		return nil
	}

	Args := info.Args()
	args_str := []string{}
	for _, v := range Args {
		args_str = append(args_str, v.String())
	}
	switch args_str[0] {
	case "require":
		v := args_str[1]
		javascript_API_lib_lock.RLock()
		module, ok := javascript_API_lib[v]
		javascript_API_lib_lock.RUnlock()

		if ok {
			var body string
			if module.enabled { // some module such as networking is disabled by default
				if !module.is_in_ram {
					b, err := ioutil.ReadFile(module.model_path)
					if err != nil {
						errs <- err
						return nil
					}
					body = string(b)
				} else {
					body = module.script
				}

				val, _ := v8go.NewValue(vm, string(body))

				return val
			}

		} else { // module not exist
			errs <- errors.New("could not import" + v + " (cannot find package " + v + " in any of path)")
		}

	case "Get":
		dtype, p := datastore.Get(args_str[1])
		if p == nil {
			return nil
		}
		var val *v8go.Value
		var err error
		switch dtype {
		case types.String:
			val, err = v8go.NewValue(vm, *(*string)(p))
		case types.Int, types.Int64:
			val, err = v8go.NewValue(vm, *(*int64)(p))
		case types.Int32:
			val, err = v8go.NewValue(vm, *(*int32)(p))

		case types.Int16:
			val, err = v8go.NewValue(vm, *(*int16)(p))

		case types.Int8:
			val, err = v8go.NewValue(vm, *(*int8)(p))

		case types.Int128:
			val, err = ctx.RunScript("BigInt("+(*(*int128.Int128)(p)).String()+")", "int128.js")

		}
		return checkerr(err, val, errs)

	case "Set":
	// set function will generate a reverse function
	// reverse function will be executed if any error occours

	case "Move":

	case "Copy":

	case "Find":

	case "create":
		val, err := creater(ctx, tmp_store, args_str[1:]...)
		return checkerr(err, val, errs)

	case "value":
		switch args_str[1] {
		case "call": // call calls the value as a function, this is created by the "new GoFunction()"

			fn, ok := tmp_store[args_str[2]].(reflect.Value)
			if !ok {
				errs <- errors.New("cannot call non function value")
			}
			fn_t := fn.Type()
			if fn_t.NumIn() > len(args_str[3:]) {
				errs <- errors.New("not enough argument calling function " + fn_t.Name())
			}
			args_r := stringToReflectArgs(fn_t, errs, tmp_store, args_str[3:])
			if args_r == nil {
				return nil
			}
			outputs := fn.Call(args_r)
			r := []interface{}{}
			for _, v := range outputs {
				r = append(r, v.Interface())
			}

		case "method": // callfn calls a function from type
			// call_go_fn("value", "method","pathToValue", "methodName",...args)
			ptr, ok := tmp_store[args_str[2]]
			if !ok {
				errs <- errors.New("cannot find value by path " + args_str[2])
				return nil
			}
			val := reflect.ValueOf(ptr)
			fn := val.MethodByName(args_str[3])
			if !fn.IsValid() {
				errs <- errors.New("type " + val.Type().Name() + " has no method " + args_str[3])
				return nil
			}
			args_r := stringToReflectArgs(fn.Type(), errs, tmp_store, args_str[3:])
			if args_r == nil {
				return nil
			}
			outputs := fn.Call(args_r)
			r := []interface{}{}
			for _, v := range outputs {
				r = append(r, v.Interface())
			}
		case "value":
		case "string":

		}

	case "settings":
		if uid == 19890604 {
			errs <- errors.New("permission denied")
			return nil
		}
		a := []string{}
		for _, v := range Args[1:] {
			a = append(a, v.String())
		}
		val, err := config.JsHandle(ctx, uid, gid, a...)
		return checkerr(err, val, errs)

	case "tensorflow":

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

func stringToReflectArgs(fn_t reflect.Type, errs chan error, tmp_store map[string]interface{}, args_str []string) []reflect.Value {
	var reflect_args = []reflect.Value{}

	for i := 0; i < len(args_str); i++ {

		inV := fn_t.In(i)
		in_Kind := inV.Kind() //func

		var val string

		switch kind := in_Kind.String(); kind {
		case "string":
			reflect_args = append(reflect_args, reflect.ValueOf(val))

		case "ptr":
			value, ok := tmp_store[val]
			if !ok {
				errs <- errors.New("invalid memory address, expected Go pointer values")
			}
			v := reflect.ValueOf(value)
			v_k := v.Kind()
			if v_k != reflect.Ptr {
				errs <- errors.New("expected pointer value got " + v_k.String())
				return nil
			}
			if a, b := reflect.Indirect(v).Kind(), inV.Elem().Kind(); a != b {
				errs <- errors.New("expected type " + b.String() + " got " + a.String())
				return nil
			}
			reflect_args = append(reflect_args, v)

		case "bool":
			b, err := strconv.ParseBool(val)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(b))

		case "int":
			b, err := strconv.Atoi(val)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(b))
		case "int8":
			b, err := strconv.ParseInt(val, 10, 8)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}

			}
			reflect_args = append(reflect_args, reflect.ValueOf(int8(b)))
		case "int16":
			b, err := strconv.ParseInt(val, 10, 16)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(int16(b)))
		case "int32":
			b, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(int32(b)))
		case "int64":
			b, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(int64(b)))
		case "uint":
			b, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(uint(b)))
		case "uint8":
			b, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(uint8(b)))
		case "uint16":
			b, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(uint16(b)))
		case "uint32":
			b, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(uint32(b)))
		case "uint64":
			b, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(b))
		case "uintptr":
		case "float34":
			b, err := strconv.ParseFloat(val, 32)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(float32(b)))
		case "float64":
			b, err := strconv.ParseFloat(val, 64)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(b))
		case "complex64":
			b, err := strconv.ParseComplex(val, 64)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(complex64(b)))
		case "complex128":
			b, err := strconv.ParseComplex(val, 128)
			if err != nil {
				if v, ok := tmp_store[val]; ok {
					a := reflect.ValueOf(v)
					if a.Elem().Kind().String() == kind {
						reflect_args = append(reflect_args, a.Elem())
						continue
					}
				} else {
					errs <- err
					return nil
				}
			}
			reflect_args = append(reflect_args, reflect.ValueOf(b))

		case "array", "unsafepointer", "chan", "func", "interface", "map", "slice", "struct": // a path to tmp_store
			value, ok := tmp_store[val]
			if !ok {
				errs <- errors.New("invalid memory address, expected Go pointer values")
			}
			if v := reflect.TypeOf(value).Elem().Kind().String(); v != kind {
				errs <- errors.New("expected type " + kind + " got " + v)
			}
			reflect_args = append(reflect_args, reflect.ValueOf(value).Elem())
		}

	}
}
