package router

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func write_json(data map[string]interface{}) error {

	a, err := json.Marshal(data)
	go CheckErr(err)

	err = ioutil.WriteFile("file.json", a, 0777)
	go CheckErr(err)

	return nil
}

func read_json(filename string) (map[string]interface{}, error) {
	f, err := os.Open(filename)
	go CheckErr(err)

	data, err := ioutil.ReadAll(f)

	go CheckErr(err)

	var value map[string]interface{}

	err = json.Unmarshal(data, &value)
	go CheckErr(err)

	return value, nil
}

func retreive_data() {

}
