package router

import (
	"log"
	"net"
)

// a common go file to get information
func getMacAddrs() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func get_first_mac_addr() string {
	a, _ := getMacAddrs()
	return a[0]
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getType(obj interface{}) string {
	switch obj.(type) {

	case int:
		return "int"
	case string:
		return "string"
	case map[string]interface{}:
		return "map[string]interface{}"
	case bool:
		return "bool"
	case uint:
		return "uint"

	default:
		return ""
	}
}
