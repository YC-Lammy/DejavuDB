package clustermanager

import (
	"../../register"

	"time"
)

func init() {
	go func() {
		time.Sleep(100 * time.Millisecond)
		for _, v := range register.Heartbeat_conns {

		}
	}()
}
