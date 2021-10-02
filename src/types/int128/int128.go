package int128

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <float.h>

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

__int128_t Sub(__int128_t v, long long e){
    return v-e;
}

__int128_t Add(__int128_t v, long long e){
    return v+e;
}
*/
import "C"
import "unsafe"

type Int128 [16]byte

func (i Int128) String() string {
	a := C.CString("")
	defer C.free(unsafe.Pointer(a))
	b := C.itoa128(a, [16]byte(i), 10)
	c := C.GoString(b)
	return c
}

func StrToInt128(str string) (Int128, error) {
	a := C.CString(str)
	defer C.free(unsafe.Pointer(a))
	b := C.atoi128(a)
	return Int128(b), nil
}

func (i *Int128) Add(e int) {

}

func (i *Int128) Sub(e int) {
	cint := C.longlong(e)
	*i = Int128(C.Sub([16]byte(*i), cint))
}
