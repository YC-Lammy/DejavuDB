package keyvalue

import (
	"src/types"
	"sync"
	"unsafe"
)

type Value struct {
	Dtype byte
	Value unsafe.Pointer
}

func NewValue(data interface{}) (Value, error) {
	val := Value{}
	switch v := data.(type) {
	case string:
		val.Dtype = types.String
		val.Value = unsafe.Pointer(&v)
	}
	return val, nil
}

type KeyStore struct {
	Data map[string]Value
	Lock sync.RWMutex
}

func (k *KeyStore) Add(key string, val Value) {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	k.Data[key] = val
}

func (k *KeyStore) Delete(key string) {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	delete(k.Data, key)
}
