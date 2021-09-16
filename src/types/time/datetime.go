package time

import "strconv"

type Datetime struct {
	Year       uint8
	Month      uint8
	Day        uint8
	Hour       uint8
	Minute     uint8
	Second     uint8
	Nanosecond uint16
}

type Smalldatetime uint32

func (s Smalldatetime) String() string {
	a := strconv.Itoa(int(s))

	return "20" + a[0:2] + ":" + a[2:4] + ":" + a[4:6] + " " + a[6:8] + ":" + a[8:10]
}
