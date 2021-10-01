package datastore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"src/types"
)

func init() {
	s, _ := os.UserHomeDir()

	origin := path.Join(s, "dejavuDB", "database")
	arr, _ := ioutil.ReadDir(origin)
	for _, v := range arr {
		f, err := os.Open(path.Join(origin, v.Name()))
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		v, _ := ioutil.ReadAll(f)
		ptr, dtype, err := types.FromBytes(v)
		JsSet(f.Name(), ptr, dtype)
	}
}
