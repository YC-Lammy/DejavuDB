package binjson

import (
	cbor "github.com/fxamacker/cbor/v2"
	json "github.com/goccy/go-json"
)

type BinaryJson struct {
	B []byte
} // binjson is in cbor format

func NewBinaryJson(str []byte) (*BinaryJson, error) {
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		return nil, err
	}
	b, err := cbor.Marshal(m)
	if err != nil {
		return nil, err
	}
	a := &BinaryJson{B: b}
	return a, err
}

func (b *BinaryJson) String() string {
	m := map[string]interface{}{}
	cbor.Unmarshal(b.B, &m)
	d, _ := json.Marshal(m)
	return string(d)
}

func (b *BinaryJson) Set(key string) {

}
