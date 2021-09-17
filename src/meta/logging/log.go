package logging

import (
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var log_logfile *os.File

var home_dir, _ = os.UserHomeDir()

func logHandler(message string) {
	log_logfile.WriteString(message)
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
		log.SetOutput(log_logfile)

	}

	for {
		time.Sleep(24 * time.Hour)

		d := time.Now().Day()
		if day != int(d) {
			day = int(d)

			p := strings.Join([]string{home_dir, "dejavuDB", "log", time.Now().Format("2006-01-02")}, string(os.PathSeparator))
			new, _ := os.Create(p)

			log_logfile.Close()
			log_logfile = new
			log.SetOutput(log_logfile)

		}

	}
}
