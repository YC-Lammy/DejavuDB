package datastore

import (
	"io/fs"
	"math/rand"
	"strings"
	"time"
)

type Temp struct {
	closed bool
	name   string
	data   interface{}
	Dtype  string
}

var temp_store = map[string]*Temp{}

func TempFile(dir, pattern string) (*Temp, error) {
	tmp := Temp{closed: false, name: pattern + RandStringBytesMaskImprSrcSB((10))}
	return &tmp, nil
}

func (tmp *Temp) Close() error {
	tmp.closed = true
	delete(temp_store, tmp.name)
	return nil
}

func (tmp *Temp) Name() string {
	if tmp.closed {
		return ""
	}
	return tmp.name
}

func (tmp *Temp) Read() (interface{}, error) {
	if tmp.closed {
		return nil, fs.ErrClosed
	}
	return tmp.data, nil
}

func (tmp *Temp) Write(data interface{}) {
	tmp.data = data
}

func RandStringBytesMaskImprSrcSB(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits := 6                    // 6 bits to represent a letter index
	letterIdxMask := 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax := 63 / letterIdxBits    // # of letter indices fitting in 63 bits
	src := rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(int(cache) & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
