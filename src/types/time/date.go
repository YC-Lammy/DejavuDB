package time

import (
	"errors"
	"strconv"
	"strings"
)

type Date struct { // 3 bytes
	Year  uint8 // this system could not last a thousand years
	Month uint8
	Day   uint8
}

func NewDate(s string) (*Date, error) { // yyyy-mm-dd

	s = s[1:]

	a := strings.Split(s, "-")

	if len(a) != 3 || len(s) != 9 {
		return nil, errors.New("error formating date")
	}
	y, err := strconv.Atoi(a[0])
	if err != nil {
		return nil, err
	}
	m, err := strconv.Atoi(a[1])
	if err != nil {
		return nil, err
	}
	d, err := strconv.Atoi(a[2])
	if err != nil {
		return nil, err
	}

	return &Date{uint8(y), uint8(m), uint8(d)}, nil
}
