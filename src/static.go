package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <float.h>

char *str2LongDouble (char* str)
{
	long double f1;
	f1 = strtold (str, NULL);
	char * arr= malloc(sizeof(f1));
	memcpy(arr, (char*)(&f1),sizeof(f1));
	return arr;
}

char* LongDouble2str (char* arr)
{
	long double f1;
	memcpy(&f1, arr, strlen(arr)+1);
	char s[100];
	sprintf(s,"%.*Lf",LDBL_DECIMAL_DIG,f1);
	char *str = malloc(strlen(s)+1);
	memcpy(str, s, strlen(s)+1);

	return str;
}

__int128_t atoi128(const char *s)
{
    while (*s == ' ' || *s == '\t' || *s == '\n' || *s == '+') ++s;
    int sign = 1;
    if (*s == '-')
    {
        ++s;
        sign = -1;
    }
    size_t digits = 0;
    while (s[digits] >= '0' && s[digits] <= '9') ++digits;
    char scratch[digits];
    for (size_t i = 0; i < digits; ++i) scratch[i] = s[i] - '0';
    size_t scanstart = 0;

    __int128_t result = 0;
    __int128_t mask = 1;
    while (scanstart < digits)
    {
        if (scratch[digits-1] & 1) result |= mask;
        mask <<= 1;
        for (size_t i = digits-1; i > scanstart; --i)
        {
            scratch[i] >>= 1;
            if (scratch[i-1] & 1) scratch[i] |= 8;
        }
        scratch[scanstart] >>= 1;
        while (scanstart < digits && !scratch[scanstart]) ++scanstart;
        for (size_t i = scanstart; i < digits; ++i)
        {
            if (scratch[i] > 7) scratch[i] -= 3;
        }
    }

    return result * sign;
}


char *utoa128(char *dest, __uint128_t v, int base) {
    char buf[129];
    char *p = buf + 128;
    const char *digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ";

    *p = '\0';
    if (base >= 2 && base <= 36) {
        while (v > (unsigned)base - 1) {
            *--p = digits[v % base];
            v /= base;
        }
        *--p = digits[v];
    }
    return strcpy(dest, p);
}

char *itoa128(char *buf, __int128_t v, int base) {
    char *p = buf;
    __uint128_t uv = (__uint128_t)v;
    if (v < 0) {
        *p++ = '-';
        uv = -uv;
    }
    if (base == 10)
        utoa128(p, uv, 10);
    else
    if (base == 16)
        utoa128(p, uv, 16);
    else
        utoa128(p, uv, base);
    return buf;
}


*/
import "C"
import (
	"math/big"
	"strconv"
	"time"
	"unsafe"
)

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

type float128 []byte

type bigfloat big.Float

type int128 [16]byte

type int256 struct {
	p1 int64
	p2 int64
	p3 int64
	p4 int64
}

type smallint uint16

type tinyint uint8

type bigint big.Int

type Time time.Time

type duration time.Duration

func (l float128) String() string {
	str := C.CString(string(l))
	new := C.LongDouble2str(str)
	result := C.GoString(new)

	C.free(unsafe.Pointer(str))
	C.free(unsafe.Pointer(new))
	return result
}

func strToFloat128(str string) (float128, error) {
	a := C.CString(str)
	b := C.str2LongDouble(a)
	c := C.GoString(b)
	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	return float128([]byte(c)), nil
}

func (i int128) String() string {
	a := C.CString("")
	b := C.itoa128(a, [16]byte(i), 10)
	c := C.GoString(b)
	C.free(unsafe.Pointer(a))
	return c
}

func strToInt128(str string) (int128, error) {
	a := C.CString(str)
	b := C.atoi128(a)
	C.free(unsafe.Pointer(a))
	return int128(b), nil
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
