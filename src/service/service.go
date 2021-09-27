package service

import (
	"time"
)

var Services = map[string]*Service{}

type Service struct {
	Name        string
	Type        string
	Active      bool
	ActiveSince time.Time
	Task        []string
	Cfg         string
}

func init() {
	wf := Service{
		Name:        "workflow",
		Type:        "internal",
		Active:      true,
		ActiveSince: time.Now(),
		Task:        []string{},
		Cfg:         "",
	}
	Services["workflow"] = &wf
}
