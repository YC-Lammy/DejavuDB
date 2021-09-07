package main

import (
	"time"
)

func Difference_str_arr(a, b []string) []string {

	ok := true
	var list []string

	for _, v := range b {
		for _, x := range a {
			if v == x {
				ok = false
				break
			}
		}
		if ok {
			list = append(list, v)
		}
		ok = true
	}
	return list
}

func stringSliceIndex(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}
	return -1
}

func getMacFromIp(ipv4 string) *string {
	for key, v := range shard_map {
		if v != nil {
			if v.RemoteAddr().String() == ipv4 {
				return &key
			}
		}
	}
	for key, v := range router_map {
		if v != nil {
			if v.RemoteAddr().String() == ipv4 {
				return &key
			}
		}

	}
	return nil
}

func getIpFromMac(mac string) *string {
	if v, ok := shard_map[mac]; ok {
		if v != nil {
			ip := v.RemoteAddr().String()
			return &ip
		}
	}
	if v, ok := router_map[mac]; ok {
		if v != nil {
			ip := v.RemoteAddr().String()
			return &ip
		}
	}
	return nil
}

func waitUntil(condition func() bool, execute func(), duration time.Duration) {
	for {
		if condition() {
			execute()
			break
		}
		time.Sleep(duration)
	}
}
