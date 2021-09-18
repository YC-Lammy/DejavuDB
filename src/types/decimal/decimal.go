package decimal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Decimal32 struct {
	int int16
	p   uint16
}
type Decimal64 struct {
	int int32
	p1  uint32
}

type Decimal128 struct {
	int int32
	p1  int32
	p2  int64
}

type Decimal192 struct {
	int int64
	p1  int64
	p2  int64
}

func (f Decimal32) String() string {
	a := strconv.Itoa(int(f.int))
	b := strconv.FormatUint(uint64(f.p), 10)
	return a + "." + b
}

func StrToDecimal32(str string) (Decimal32, error) {
	s := strings.Split(str, ".")
	if len(s) > 2 {
		return Decimal32{}, errors.New("cannot convert " + str + " to float base 10")
	} else if len(s) == 1 {
		i, err := strconv.ParseUint(s[0], 10, 16)
		if err != nil {
			return Decimal32{}, err
		}
		if str[0] == '.' {

			return Decimal32{int: int16(0), p: uint16(i)}, nil
		} else {
			return Decimal32{int: int16(i), p: uint16(0)}, nil
		}
	}
	i, err := strconv.ParseInt(s[0], 10, 16)
	if err != nil {
		return Decimal32{}, err
	}
	p1, err := strconv.ParseUint(s[1], 10, 16)
	if err != nil {
		return Decimal32{}, errors.New(strings.Replace(err.Error(), "strconv.ParseInt: ", "", -1) + "try using ds128")
	}
	return Decimal32{int: int16(i), p: uint16(p1)}, nil
}

func (f Decimal64) String() string {
	a := strconv.Itoa(int(f.int))
	b := strconv.FormatUint(uint64(f.p1), 10)
	return a + "." + b
}

func StrToDecimal64(str string) (Decimal64, error) {
	s := strings.Split(str, ".")
	if len(s) > 2 {
		return Decimal64{}, errors.New("cannot convert " + str + " to float base 10")
	} else if len(s) == 1 {
		i, err := strconv.ParseUint(s[0], 10, 32)
		if err != nil {
			return Decimal64{}, err
		}
		if str[0] == '.' {

			return Decimal64{int: int32(0), p1: uint32(i)}, nil
		} else {
			return Decimal64{int: int32(i), p1: uint32(0)}, nil
		}
	}
	i, err := strconv.ParseInt(s[0], 10, 32)
	if err != nil {
		return Decimal64{}, err
	}
	p1, err := strconv.ParseUint(s[1], 10, 32)
	if err != nil {
		return Decimal64{}, errors.New(strings.Replace(err.Error(), "strconv.ParseInt: ", "", -1) + "try using ds128")
	}
	return Decimal64{int: int32(i), p1: uint32(p1)}, nil
}

func (f Decimal128) String() string {
	return fmt.Sprintf("%v.%v%v", f.int, f.p1, f.p2)
}

func StrToDecimal128(str string) (Decimal128, error) {
	s := strings.Split(str, ".")
	if len(s) > 2 {
		return Decimal128{}, errors.New("cannot convert " + str + " to float base 10")
	}

	var i int64
	if str[0] == '.' {
		i = 0
	} else {
		a, err := strconv.ParseInt(s[0], 10, 32)
		if err != nil {
			return Decimal128{}, err
		}
		i = a
	}

	var f1 string
	var f2 string

	if len(s) == 2 {
		f1 = s[1][:len(s[1])/3]
		f2 = s[1][len(s[1])/3:]
	} else if len(s) == 1 {
		f1 = s[0][:len(s[1])/3]
		f2 = s[0][len(s[1])/3:]
	} else {
		return Decimal128{}, errors.New("cannot convert " + str + " to float base 10")
	}

	p1, err := strconv.ParseInt(f1, 10, 32)
	if err != nil {
		return Decimal128{}, err
	}
	p2, err := strconv.ParseInt(f2, 10, 64)
	if err != nil {
		return Decimal128{}, err
	}
	return Decimal128{int: int32(i), p1: int32(p1), p2: p2}, nil

}

func (f Decimal192) String() string {
	return fmt.Sprintf("%v.%v%v", f.int, f.p1, f.p2)
}

func StrToDecimal192(str string) (Decimal192, error) {
	s := strings.Split(str, ".")
	if len(s) > 2 {
		return Decimal192{}, errors.New("cannot convert " + str + " to float base 10")
	}

	var i int64
	if str[0] == '.' {
		i = 0
	} else {
		a, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			return Decimal192{}, err
		}
		i = a
	}

	var f1 string
	var f2 string

	if len(s) == 2 {
		f1 = s[1][:len(s[1])/2]
		f2 = s[1][len(s[1])/2:]
	} else if len(s) == 1 {
		f1 = s[0][:len(s[1])/2]
		f2 = s[0][len(s[1])/2:]
	} else {
		return Decimal192{}, errors.New("cannot convert " + str + " to float base 10")
	}

	p1, err := strconv.ParseInt(f1, 10, 64)
	if err != nil {
		return Decimal192{}, err
	}
	p2, err := strconv.ParseInt(f2, 10, 64)
	if err != nil {
		return Decimal192{}, err
	}
	return Decimal192{int: int64(i), p1: p1, p2: p2}, nil

}
