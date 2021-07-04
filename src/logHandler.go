package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var log_logfile *os.File

func logHandler(conn net.Conn, message string) {
	log.SetOutput(log_logfile)
	switch strings.Split(message, " ")[0] {

	}
	log.Println(message)
}

func log_file_date() { // this func loop once every day and create a new log file at 00:00:00

	var day int

	log_logfile, _ = os.Create(time.Now().Format("2006-01-02"))

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

		new, err := os.Create(time.Now().Format("2006-01-02"))
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

			new, err := os.Create(time.Now().Format("2006-01-02"))
			if err != nil {
				log.Println(err)
				return
			}

			log_logfile.Close()
			log_logfile = new

		}

	}
}
