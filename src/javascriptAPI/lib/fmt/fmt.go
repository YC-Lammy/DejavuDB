package fmtjs

import "fmt"

func JsHandle(args ...string) string {
	switch args[0] {
	case "Print":
		a := []interface{}{}
		for _, v := range args[1:] {
			a = append(a, v)
		}
		return fmt.Sprint(a...)
	case "Printf":
		a := []interface{}{}
		for _, v := range args[2:] {
			a = append(a, v)
		}
		return fmt.Sprintf(args[1], a...)
	case "Println":
		a := []interface{}{}
		for _, v := range args[1:] {
			a = append(a, v)
		}
		return fmt.Sprintln(a...)

	}
	return ""
}
