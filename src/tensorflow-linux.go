// +build linux

package main

import "runtime"

func init_tensorflow() error {
	switch runtime.GOOS {
	case "linux":
	case "darwin":
	case "windows":

	}
	return nil
}
