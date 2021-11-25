package logging

import (
	"dejavuDB/src/config"
	"io"
	"log"
	"os"
	"path"
	"time"
)

var _ io.Writer = (*FileLogger_t)(nil)

var Logger = NewFileLogger()

type FileLogger_t struct {
	*log.Logger
	DefaultWriter io.Writer
	FileWriter    io.Writer

	nextInitialTime time.Time
}

func NewFileLogger() *FileLogger_t {
	t := time.Now()
	p := path.Join(config.RootDir, "log")
	f, _ := os.Create(path.Join(p, t.Format("2006-01-02")))

	l := &FileLogger_t{
		Logger:          log.Default(),
		nextInitialTime: time.Date(t.Year(), t.Month(), t.Day()+1, 00, 0, 0, 0, t.Location()),
		DefaultWriter:   log.Writer(),
		FileWriter:      f,
	}

	l.Logger.SetOutput(l)

	return l
}

func (l *FileLogger_t) Write(p []byte) (n int, err error) {

	n, err = l.DefaultWriter.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		err = io.ErrShortWrite
		return
	}

	if time.Now().Sub(l.nextInitialTime) >= 0 {
		p := path.Join(config.RootDir, "log")
		f, er := os.Create(path.Join(p, time.Now().Format("2006-01-02")))
		if err != nil {
			err = er
			return
		}
		l.FileWriter = f
	}
	n, err = l.FileWriter.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		err = io.ErrShortWrite
		return
	}
	return len(p), nil
}
