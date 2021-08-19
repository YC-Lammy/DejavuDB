package tensorflow

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func init_python_server() error {
	switch runtime.GOOS {

	case "linux":
		script := `
		#!/bin/bash

		for version in "python3.9" "python3.8" "python3.7" "python3.6" "python3.5" "python2.7" "python2.6" 
		do
			if type $version &> /dev/null; then
				echo "found" $version
				exit 0
			fi
		done
		echo not found`

		tmpFile, err := ioutil.TempFile(os.TempDir(), "sh-*.sh")
		if err != nil {
			panic(err)
		}
		tmpFile.WriteString(script)
		cmd := exec.Command("bash", tmpFile.Name())
		out, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		fmt.Println(string(out))

		if strings.Replace(string(out), "\n", "", -1) == "not found" {
			fmt.Println("python not installed (version 2.6, 2.7, 3.5+)")
			os.Exit(1)
		}
		tmpFile.Close()

	case "darwin":
	case "windows":

	}
	return nil
}

var python_server_script = `
# copy from tensorflow/tenserflowHost.py
`
