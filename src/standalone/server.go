package standalone

import (
	"io"
	"log"
	"net"
	"os"
	"path"
	"src/config"
	"src/standalone/client_interface"
	"time"
)

var log_logfile *os.File

func Start() {
	multiWriter := io.MultiWriter(log.Writer(), log_logfile)
	log.SetOutput(multiWriter)
	c, err := net.Listen("tcp", config.Client_port) // client interface
	if err != nil {
		panic(err)
	}
	defer c.Close()
	for {
		conn, err := c.Accept()
		if err != nil {
			log.Println(err)
		}
		go client_interface.Handle(conn)
	}
}

func log_file_date() { // this func loop once every day and create a new log file at 00:00:00

	var day int

	home_dir, _ := os.UserHomeDir()

	log_path := path.Join(home_dir, "dejavuDB", "log")

	f, _ := os.Create(path.Join(log_path, time.Now().Format("2006-01-02")))
	*log_logfile = *f

	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day()+1, 00, 0, 0, 1, t.Location())
	duration := n.Sub(t)
	if duration < 0 {
		n = n.Add(24 * time.Hour)
		duration = n.Sub(t)
	}

	time.Sleep(duration) // sleep through the first day until 00:00
	d := time.Now().Day()
	if day != int(d) {
		day = int(d)

		new, err := os.Create(path.Join(log_path, time.Now().Format("2006-01-02")))
		if err != nil {
			log.Println(err)
			return
		}

		log_logfile.Close()
		*log_logfile = *new

	}

	for {
		time.Sleep(24 * time.Hour)

		d := time.Now().Day()
		if day != int(d) {
			day = int(d)

			new, _ := os.Create(path.Join(log_path, time.Now().Format("2006-01-02")))

			log_logfile.Close()
			*log_logfile = *new

		}

	}
}
