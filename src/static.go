package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

__uint128_t str2LongDouble (char* str)
{
	long double f1;
	__int128_t arr;
	f1 = strtold (str, NULL);
	memcpy(&arr,(unsigned char*)(&f1),sizeof(f1));
	return arr;
}
char* LongDouble2str (__uint128_t arr)
{
	long double f1;
	char* str;
	memcpy(&f1,(unsigned char*)(&arr),sizeof(f1));
	sprintf(str,"%Lf",f1);
	return str;
}
__uint128_t str2Int128 (char* str)
{
	__int128_t arr;
	return arr;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type float128 struct {
	int int
	p1  int64
	p2  int64
}

type int128 struct {
	p1 int64
	p2 int64
}

type int256 struct {
	p1 int64
	p2 int64
	p3 int64
	p4 int64
}

func (f float128) String() string {
	return fmt.Sprintf("%v.%v%v", f.int, f.p1, f.p2)
}

func strToFloat128(str string) (float128, error) {
	s := strings.Split(str, ".")
	if len(s) != 2 {
		return float128{}, errors.New("cannot convert " + str + " to float base 10")
	}
	f1 := s[1][:len(s[1])/2]
	f2 := s[1][len(s[1])/2:]
	i, err := strconv.ParseInt(s[0], 10, 64)
	if err != nil {
		return float128{}, err
	}
	p1, err := strconv.ParseInt(f1, 10, 64)
	if err != nil {
		return float128{}, err
	}
	p2, err := strconv.ParseInt(f2, 10, 64)
	if err != nil {
		return float128{}, err
	}
	return float128{int: int(i), p1: p1, p2: p2}, nil

}

func (i int128) String() string {
	return strconv.FormatInt(i.p1, 10) + strconv.FormatInt(i.p2, 10)
}

func strToInt128(str string) (int128, error) {
	p1 := str[:len(str)/2]
	p2 := str[len(str)/2:]
	a, err := strconv.ParseInt(p1, 10, 64)
	if err != nil {
		return int128{}, err
	}
	b, err := strconv.ParseInt(p2, 10, 64)
	if err != nil {
		return int128{}, err
	}
	return int128{p1: a, p2: b}, nil
}

func (i int256) String() string {
	return strconv.FormatInt(i.p1, 10) + strconv.FormatInt(i.p2, 10) + strconv.FormatInt(i.p3, 10) + strconv.FormatInt(i.p4, 10)
}

func strToInt256(str string) (int256, error) {
	i := len(str) / 4
	p1 := str[:i]
	p2 := str[i : len(str)/2]
	p3 := str[len(str)/2 : i*3]
	p4 := str[i*3:]

	a, err := strconv.ParseInt(p1, 10, 64)
	if err != nil {
		return int256{}, err
	}
	b, err := strconv.ParseInt(p2, 10, 64)
	if err != nil {
		return int256{}, err
	}
	c, err := strconv.ParseInt(p3, 10, 64)
	if err != nil {
		return int256{}, err
	}
	d, err := strconv.ParseInt(p4, 10, 64)
	if err != nil {
		return int256{}, err
	}
	return int256{p1: a, p2: b, p3: c, p4: d}, nil
}
