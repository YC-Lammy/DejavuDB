package workflow

import "unsafe"

const (
	duration_method = "duration"
	date_method     = "date"
	weekday_method  = "weekday"
	email_method    = "email"
	event_method    = "event"
)

const (
	all_success            = 0x00
	all_failed             = 0x01
	all_done               = 0x02
	one_failed             = 0x03
	one_sucess             = 0x04
	none_failed            = 0x05
	none_failed_or_skipped = 0x06
	none_skipped           = 0x07
	dummy                  = 0x08
)

func Chain(works ...*Work) {
	for i := 0; i < len(works)-1; i++ {
		works[i].SetDownstream(works[i+1])
	}
}

type Workflow struct {
	Start *Work

	Trigger_method string
	Trigger        unsafe.Pointer

	Email_from     string
	Email_password string
	Email_to       []string
	Email_stmp     string // stmp server host and port
}

type Work struct {
	Id string

	Trigger_rule byte

	Dependencies []*Work
	Retry        uint8
	Timeout      uint64 //ms, if 0 no timeout

	Script string
}

func (w *Work) SetDownstream(works ...*Work) {}

func (w *Work) SetUpstream(works ...*Work) {

}
