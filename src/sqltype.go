package main

/*
# include <string.h>
# include <stdlib.h>

char * intToByte(long long value){
	char * str = malloc(sizeof(value));
    memcpy(str, &value, sizeof(value));
	return str;
}
char * floatToByte(double value){
	char * str = malloc(sizeof(value));
    memcpy(str, &value, sizeof(value));
	return str;
}
*/
import "C"
import (
	"io"
	"unsafe"
)

type BLOB struct {
	data     []byte
	readpos  int
	readdone bool
}

func (blob *BLOB) Write(values ...interface{}) (int, error) {
	buffer := []byte{}
	for _, v := range values {
		switch v := v.(type) {
		case string:
			buffer = append(buffer, []byte(v)...)
		case []byte:
			buffer = append(buffer, []byte(v)...)
		case byte:
			buffer = append(buffer, v)
		case int:
			buffer = append(buffer, intToByte(v)...)
		case int64:
			buffer = append(buffer, intToByte(int(v))...)
		case int32:
			buffer = append(buffer, intToByte(int(v))...)
		case int16:
			buffer = append(buffer, intToByte(int(v))...)
		case int8:
			buffer = append(buffer, byte(v))
		case bool:
			if v {
				buffer = append(buffer, byte(1))
			} else {
				buffer = append(buffer, byte(0))
			}
		case float32:
			buffer = append(buffer, floatToByte(float64(v))...)
		case float64:
			buffer = append(buffer, floatToByte(v)...)
		}
	}
	blob.data = append(blob.data, buffer...)
	blob.readdone = false
	return len(buffer), nil
}

func (blob *BLOB) Update(values ...interface{}) (int, error) {

	blob.data = []byte{}
	return blob.Write(values...)
}

func (blob *BLOB) Read(p []byte) (int, error) {
	if blob.readdone {
		return 0, io.EOF
	}
	length := len(blob.data)
	i := 0
	for i < len(p) {
		if i+blob.readpos <= length {
			p[i] = blob.data[i+blob.readpos]
		} else {
			blob.readdone = true
			break
		}
		i++
	}
	return i, nil
}

func (blob *BLOB) ReadAll() ([]byte, error) {
	return blob.data, nil
}

func intToByte(value int) []byte {
	a := C.longlong(value)
	b := C.intToByte(a)
	c := C.GoString(b)
	C.free(unsafe.Pointer(&a))
	C.free(unsafe.Pointer(b))
	return []byte(c)
}
func floatToByte(value float64) []byte {
	a := C.double(value)
	b := C.floatToByte(a)
	c := C.GoString(b)
	C.free(unsafe.Pointer(&a))
	C.free(unsafe.Pointer(b))
	return []byte(c)
}
