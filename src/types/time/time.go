package time

import (
	"errors"
	"strconv"
	"strings"
)

type Time struct { // 5 bytes
	Hour       uint8
	Minute     uint8
	Second     uint8
	Nanosecond uint16
}

func NewTime(h uint8, m uint8, s uint8, n uint16) (*Time, error) {
	return &Time{Hour: h, Minute: m, Second: s, Nanosecond: n}, nil
}

func ParseTime(s string) (*Time, error) {
	var err error
	var h int
	var m int
	var sec int
	sep := strings.Split(s, ":")
	if len(sep) != 3 {
		return nil, errors.New("time missing field")
	}

	se := strings.Split(sep[2], ".")
	if len(se) == 0 || len(se) > 2 {
		return nil, errors.New("time missing field")
	}

	if h, err = strconv.Atoi(sep[0]); h > 24 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("hours larger than 24 interval")
	}
	if m, err = strconv.Atoi(sep[1]); m > 59 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Minutes larger than 59 interval")
	}
	if sec, err := strconv.Atoi(se[0]); sec > 59 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Minutes larger than 59 interval")
	}

	time := &Time{Hour: uint8(h), Minute: uint8(m), Second: uint8(sec)}

	if len(se) == 1 {
		time.Nanosecond = 0
	} else {
		i, err := strconv.Atoi(se[1])
		if err != nil {
			return nil, err
		}
		if i > 65535 {
			i, _ = strconv.Atoi(se[1][:len(se[1])-1])
			time.Nanosecond = uint16(i)
		} else {
			time.Nanosecond = uint16(i)
		}
	}
	return time, nil
}
