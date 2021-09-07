package datastore

import (
	"os"
)

func init(){
	
	arr, _ := ioutil.ReadDir(path.Join(os.UserHomeDir(),"dejavuDB", "database"))
	for _, v := range arr {
		var new = user{}
		f, err := os.Open(v.Name())
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		v, _ := ioutil.ReadAll(f)
	}
}