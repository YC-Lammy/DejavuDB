package float128

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
*/
import "C"
import "unsafe"

type float128 []byte

func (l float128) String() string {
	str := C.CString(string(l))
	new := C.LongDouble2str(str)
	result := C.GoString(new)

	C.free(unsafe.Pointer(str))
	C.free(unsafe.Pointer(new))
	return result
}

func StrToFloat128(str string) (float128, error) {
	a := C.CString(str)
	b := C.str2LongDouble(a)
	c := C.GoString(b)
	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	return float128([]byte(c)), nil
}
