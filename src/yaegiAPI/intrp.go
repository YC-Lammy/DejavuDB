package yaegiAPI

import (
	"bytes"

	"github.com/traefik/yaegi/interp"
)

func Run(script string) (b []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	tmp := bytes.NewBuffer([]byte{})
	i := interp.New(interp.Options{Stdout: tmp})
	v, err := i.Eval(init_script + script + ending)
	if err != nil {
		return nil, err
	}
	fn := v.Interface().(func(db) error)
	DB := &database{}

	err = fn(DB)
	if err != nil {
		return nil, err
	}

	return tmp.Bytes(), nil
}
