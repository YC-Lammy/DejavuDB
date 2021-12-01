package router_interface

import (
	"time"
)

func init() {
	go func() {
		time.Sleep(100 * time.Millisecond)
		if IsLeader {
			for k, v := range Router_conns {
				_, err := v.Heartbeat.Write([]byte{0x00})
				if err != nil {
					delete(Router_conns, k)
				}
			}
		}

	}()
}
