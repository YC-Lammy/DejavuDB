package uint128

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <float.h>

__uint128_t atou128(const char *s)
{
    while (*s == ' ' || *s == '\t' || *s == '\n' || *s == '+' || *s == '-') ++s;

    size_t digits = 0;
    while (s[digits] >= '0' && s[digits] <= '9') ++digits;
    char scratch[digits];
    for (size_t i = 0; i < digits; ++i) scratch[i] = s[i] - '0';
    size_t scanstart = 0;

    __uint128_t result = 0;
    __uint128_t mask = 1;
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

    return result;
}

char *utoa128_(char *dest, __uint128_t v, int base) {
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

char *utoa128_t(char *buf, __uint128_t uv, int base) {
    char *p = buf;
    if (base == 10)
        utoa128_(p, uv, 10);
    else
    if (base == 16)
        utoa128_(p, uv, 16);
    else
        utoa128_(p, uv, base);
    return buf;
}
*/
import "C"
import "unsafe"

type Uint128 [16]byte // a true __uint128_t

func StrToUint128(s string) (Uint128, error) {
	a := C.CString(s)
	defer C.free(unsafe.Pointer(a))
	b := C.atou128(a)

	return Uint128(b), nil
}

func (u Uint128) String() string {
	a := C.CString("")
	defer C.free(unsafe.Pointer(a))
	b := C.utoa128_t(a, [16]byte(u), 10)
	defer C.free(unsafe.Pointer(b))
	c := C.GoString(b)
	return c
}
