package datastore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"dejavuDB/src/types"
)

func init() {
	s, _ := os.UserHomeDir()

	origin := path.Join(s, "dejavuDB", "database")
	arr, _ := ioutil.ReadDir(origin)
	for _, v := range arr {
		f, err := os.Open(path.Join(origin, v.Name()))
		if err != nil {
			fmt.Println(err)
			continue
		}
		v, _ := ioutil.ReadAll(f)
		val, err := types.ValueFromBytes(v)
		if err != nil {
			continue
		}
		Set(f.Name(), val)
	}
}
