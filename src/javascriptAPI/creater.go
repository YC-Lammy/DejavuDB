package javascriptAPI

import (
	"strconv"
	"strings"

	"src/types/graph"

	"github.com/traefik/yaegi/interp"

	"rogchap.com/v8go"
)

func creater(ctx *v8go.Context, store map[string]interface{}, args ...string) (*v8go.Value, error) {
	switch strings.ToLower(args[0]) {
	case "fn":
		i := interp.New(interp.Options{})

		val, err := i.Eval(args[1])
		if err != nil {
			return nil, err
		}
		l := "path" + strconv.Itoa(len(store))
		store[l] = val
		return v8go.NewValue(vm, l)

	case "graph":
		l := "path" + strconv.Itoa(len(store))
		store[l] = graph.NewGraph()
		return v8go.NewValue(vm, l)
	case "graphvertex":
		l := "path" + strconv.Itoa(len(store))
		store[l] = graph.NewVertex()
		return v8go.NewValue(vm, l)

	case "table":
	case "json":

	case "string":
	case "int":
	case "int8":
	}
	return nil, nil
}
