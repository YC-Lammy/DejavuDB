package decimal

type decimal64 struct {
	int int32
	p1  int32
}

type decimal128 struct {
	int int32
	p1  int32
	p2  int64
}

type decimal192 struct {
	int int64
	p1  int64
	p2  int64
}

func (f decimal64) String() string {
	return fmt.Sprintf("%v.%v", f.int, f.p1)
}

func strToDecimal64(str string) (decimal64, error) {
	s := strings.Split(str, ".")
	if len(s) > 2 {
		return decimal64{}, errors.New("cannot convert " + str + " to float base 10")
	} else if len(s) == 1 {
		i, err := strconv.ParseInt(s[0], 10, 32)
		if err != nil {
			return decimal64{}, err
		}
		if str[0] == '.' {

			return decimal64{int: int32(0), p1: int32(i)}, nil
		} else {
			return decimal64{int: int32(i), p1: int32(0)}, nil
		}
	}
	i, err := strconv.ParseInt(s[0], 10, 32)
	if err != nil {
		return decimal64{}, err
	}
	p1, err := strconv.ParseInt(s[1], 10, 32)
	if err != nil {
		return decimal64{}, errors.New(strings.Replace(err.Error(), "strconv.ParseInt: ", "", -1) + "try using ds128")
	}
	return decimal64{int: int32(i), p1: int32(p1)}, nil
}

func (f decimal128) String() string {
	return fmt.Sprintf("%v.%v%v", f.int, f.p1, f.p2)
}

func strToDecimal128(str string) (decimal128, error) {
	s := strings.Split(str, ".")
	if len(s) > 2 {
		return decimal128{}, errors.New("cannot convert " + str + " to float base 10")
	}

	var i int64
	if str[0] == '.' {
		i = 0
	} else {
		a, err := strconv.ParseInt(s[0], 10, 32)
		if err != nil {
			return decimal128{}, err
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
		return decimal128{}, errors.New("cannot convert " + str + " to float base 10")
	}

	p1, err := strconv.ParseInt(f1, 10, 32)
	if err != nil {
		return decimal128{}, err
	}
	p2, err := strconv.ParseInt(f2, 10, 64)
	if err != nil {
		return decimal128{}, err
	}
	return decimal128{int: int32(i), p1: int32(p1), p2: p2}, nil

}

func (f decimal192) String() string {
	return fmt.Sprintf("%v.%v%v", f.int, f.p1, f.p2)
}

func strToDecimal192(str string) (decimal192, error) {
	s := strings.Split(str, ".")
	if len(s) > 2 {
		return decimal192{}, errors.New("cannot convert " + str + " to float base 10")
	}

	var i int64
	if str[0] == '.' {
		i = 0
	} else {
		a, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			return decimal192{}, err
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
		return decimal192{}, errors.New("cannot convert " + str + " to float base 10")
	}

	p1, err := strconv.ParseInt(f1, 10, 64)
	if err != nil {
		return decimal192{}, err
	}
	p2, err := strconv.ParseInt(f2, 10, 64)
	if err != nil {
		return decimal192{}, err
	}
	return decimal192{int: int64(i), p1: p1, p2: p2}, nil

}