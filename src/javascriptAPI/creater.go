package javascriptAPI

import (
	"strconv"
	"strings"

	"rogchap.com/v8go"

	"src/types/graph"
)

func creater(store map[string]interface{}, args ...string) (*v8go.Value, error) {
	switch strings.ToLower(args[0]) {
	case "fn":
		i := interp.New(interp.Options{})

		val, err := i.Eval(args[1])
		if err != nil {
			return nil, err
		}
		l := strconv.Itoa(len(store))
		store[l] = val
		return v8go.NewValue(vm, l)

	case "graph":
		l := strconv.Itoa(len(store))
		store[l] = graph.NewGraph()
		return v8go.NewValue(vm, l)
	case "graphvertex":
		l := uint64(len(store))
		store[l] = graph.NewVertex()
		return v8go.NewValue(vm, l)

	case "table":
	case "json":
	}
	return nil, nil
}
