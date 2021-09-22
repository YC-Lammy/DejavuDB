package javascriptAPI

import (
	"errors"
	"io/ioutil"
	"src/datastore"
	fmtjs "src/javascriptAPI/lib/fmt"
	"src/types"
	"src/types/int128"
	"strconv"

	"../settings"
	"rogchap.com/v8go"
)

func callbackfn(info *v8go.FunctionCallbackInfo, ctx *v8go.Context, errs chan error, delayfn chan *func(), args map[string]string, tmp_store map[uint64]interface{}) *v8go.Value { // when the JS function is called this Go callback will execute

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
	switch Args[0].String() {
	case "require":
		v := Args[1].String()
		module, ok := javascript_API_lib[v]

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
		dtype, p := datastore.Get(Args[1].String())
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
	// this function will be executed if any error occours

	case "create":
		val, err := creater(tmp_store, args_str[1:]...)
		return checkerr(err, val, errs)

	case "value":
		switch args_str[1] {
		case "call": // call calls the value as a function, this is created by the "new GoFunction()"
			defer func() {
				if err := recover(); err != nil {
					errs <- err.(error)
				}
			}()
		case "callfn": // callfn calls a function from type
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
		b, err := settings.JsHandle(a...)
		if err != nil {
			errs <- err
			return nil
		}
		val, err := v8go.NewValue(vm, b)
		return checkerr(err, val, errs)

	case "tensorflow":

	case "fmt":
		a := []string{}
		for _, v := range Args[1:] {
			a = append(a, v.String())
		}
		b := fmtjs.JsHandle(a...)
		val, err := v8go.NewValue(vm, b)

		return checkerr(err, val, errs)

	case "http":

		// you can return a value back to the JS caller if required
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
