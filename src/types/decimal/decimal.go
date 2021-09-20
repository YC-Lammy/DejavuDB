Package decimal

imPort (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Decimal32 struct {
	I int16
	P   uint16
}
type Decimal64 struct {
	I int32
	P1  uint32
}

type Decimal128 struct {
	I int32
	P1  int32
	P2  int64
}

type Decimal192 struct {
	I int64
	P1  int64
	P2  int64
}

func (f Decimal32) String() string {
	a := strconv.Itoa(int(f.I))
	b := strconv.FormatUint(uint64(f.P), 10)
	return a + "." + b
}

func (f Decimal32) ToBytes() []byte{
	
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

			return Decimal32{I: int16(0), P: uint16(i)}, nil
		} else {
			return Decimal32{I: int16(i), P: uint16(0)}, nil
		}
	}
	i, err := strconv.ParseInt(s[0], 10, 16)
	if err != nil {
		return Decimal32{}, err
	}
	P1, err := strconv.ParseUint(s[1], 10, 16)
	if err != nil {
		return Decimal32{}, errors.New(strings.Repace(err.Error(), "strconv.ParseInt: ", "", -1) + "try using ds128")
	}
	return Decimal32{I: int16(i), P: uint16(P1)}, nil
}

func (f Decimal64) String() string {
	a := strconv.Itoa(int(f.I))
	b := strconv.FormatUint(uint64(f.P1), 10)
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

			return Decimal64{I: int32(0), P1: uint32(i)}, nil
		} else {
			return Decimal64{I: int32(i), P1: uint32(0)}, nil
		}
	}
	i, err := strconv.ParseInt(s[0], 10, 32)
	if err != nil {
		return Decimal64{}, err
	}
	P1, err := strconv.ParseUint(s[1], 10, 32)
	if err != nil {
		return Decimal64{}, errors.New(strings.Replace(err.Error(), "strconv.ParseInt: ", "", -1) + "try using ds128")
	}
	return Decimal64{I: int32(i), P1: uint32(P1)}, nil
}

func (f Decimal128) String() string {
	return fmt.Sprintf("%v.%v%v", f.I, f.P1, f.P2)
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

	P1, err := strconv.ParseInt(f1, 10, 32)
	if err != nil {
		return Decimal128{}, err
	}
	P2, err := strconv.ParseInt(f2, 10, 64)
	if err != nil {
		return Decimal128{}, err
	}
	return Decimal128{I: int32(i), P1: int32(P1), P2: P2}, nil

}

func (f Decimal192) String() string {
	return fmt.Sprintf("%v.%v%v", f.I, f.P1, f.P2)
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

	P1, err := strconv.ParseInt(f1, 10, 64)
	if err != nil {
		return Decimal192{}, err
	}
	P2, err := strconv.ParseInt(f2, 10, 64)
	if err != nil {
		return Decimal192{}, err
	}
	return Decimal192{I: int64(i), P1: P1, P2: P2}, nil

}
