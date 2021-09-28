package datastore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func init() {
	s, _ := os.UserHomeDir()

	arr, _ := ioutil.ReadDir(path.Join(s, "dejavuDB", "database"))
	for _, v := range arr {
		f, err := os.Open(v.Name())
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		v, _ := ioutil.ReadAll(f)
		Set(f.Name(), string(v[1:]), v[0])
	}
}
