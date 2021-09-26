package tensorflow

import (
	_ "embed"
	"log"
	"net"
	"os"
	"os/exec"
	"src/config"
	"time"
)

//go:embed tensorflow.js
var tensorflowjs []byte
var Service chan []byte
var Tf_onoff chan bool

var conn net.Conn

func init() {
	go func() {
		time.Sleep(1 * time.Second)
		f, err := os.Create(os.TempDir() + string(os.PathSeparator) + "tf.js")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.Write(tensorflowjs)

		c := exec.Command("node", f.Name())
		started := false

		defer c.Process.Kill()

		if config.Enable_ML {
			c.Start()
			started = true
			con, err := net.Dial("tcp", "localhost:5630")
			if err != nil {
				log.Println(err)
			}
			conn = con
		}
		for {
			select {
			case onoff := <-Tf_onoff:
				if onoff && !started {
					c.Start()
					con, err := net.Dial("tcp", "localhost:5630")
					if err != nil {
						log.Println(err)
					}
					conn = con
				} else {
					Service <- []byte("terminate")
					c.Wait()
				}
			case cmd := <-Service:
				_, err := conn.Write(cmd)
				if err != nil {
					log.Println(err)
				}
			}
		}

	}()
}
