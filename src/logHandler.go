package main

import (
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

var log_logfile *os.File

var home_dir, _ = os.UserHomeDir()

func logHandler(conn net.Conn, message string) {
	log.SetOutput(log_logfile)
	switch strings.Split(message, " ")[0] {

	case "log":
		log_logfile.WriteString(message[3:] + "\n")
	case "getMonitor":
		send_to_all_router([]byte("monitor"))

	}
	log.Println(message)
}

func log_file_date() { // this func loop once every day and create a new log file at 00:00:00

	var day int

	os.Chdir("log")
	log_logfile, _ = os.Create(time.Now().Format("2006-01-02"))
	os.Chdir(path.Join(home_dir, "dejavuDB"))

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

		os.Chdir("log")
		new, err := os.Create(time.Now().Format("2006-01-02"))
		os.Chdir(path.Join(home_dir, "dejavuDB"))
		if err != nil {
			log.Println(err)
			return
		}

		log_logfile.Close()
		log_logfile = new

	}

	for {
		time.Sleep(24 * time.Hour)

		d := time.Now().Day()
		if day != int(d) {
			day = int(d)

			os.Chdir("log")
			new, _ := os.Create(time.Now().Format("2006-01-02"))
			os.Chdir(path.Join(home_dir, "dejavuDB"))

			log_logfile.Close()
			log_logfile = new

		}

	}
}

func sendLog(message []byte) (int, error) {
	for _, conn := range log_servers {
		n, err := send(conn, []byte("log "+string(message)))
		if err != nil {
			return n, err
		}
	}
	return len(message), nil
}

type mywrite struct {
}

func (m *mywrite) Write(p []byte) (n int, err error) {
	return sendLog(p)
}

func setupLog() {

	mywriter := &mywrite{}
	multiWriter := io.MultiWriter(log.Writer(), mywriter)
	log.SetOutput(multiWriter)
}
