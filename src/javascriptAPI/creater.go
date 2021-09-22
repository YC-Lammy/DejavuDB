package javascriptAPI

import (
	"strings"

	"rogchap.com/v8go"

	"src/types/graph"
)

func creater(store map[uint64]interface{}, args ...string) (*v8go.Value, error) {
	switch strings.ToLower(args[0]) {
	case "fn":
		i := interp.New(interp.Options{})

		val, err := i.Eval(args[1])
		if err != nil {
			return nil, err
		}
		l := uint64(len(store))
		store[l] = val
		return v8go.NewValue(vm, l)
	case "graph":
		l := uint64(len(store))
		store[l] = graph.NewGraph()
		return v8go.NewValue(vm, l)
	case "graphvertex":
		l := uint64(len(store))
		store[l] = graph.NewVertex()
		return v8go.NewValue(vm, l)
	}
	return nil, nil
}
