package yaegiAPI

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	unsf "github.com/traefik/yaegi/stdlib/unsafe"
)

func Run(script string) (b []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	tmp := bytes.NewBuffer([]byte{})
	i := interp.New(interp.Options{Stdout: tmp})
	i.Use(stdlib.Symbols)
	i.Use(unsf.Symbols)
	_, err = i.Eval(init_script + script + ending)
	if err != nil {
		return nil, err
	}
	v, err := i.Eval("foo.m")
	if err != nil {
		return nil, err
	}
	fn, ok := v.Interface().(func(interface{}) error)
	if !ok {
		return nil, errors.New("error formatting function")
	}
	//var DB db
	DB := &Database{}

	err = fn(DB)
	if err != nil {
		return nil, err
	}
	_ = fmt.Sprint(DB)

	return tmp.Bytes(), nil
}
