package fileservice

import (
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const (
	ReadAll = 0x00
	ReadAt  = 0x01
	Read    = 0x02
	Write   = 0x03
)

var opened_file map[string]*file
var map_lock = sync.RWMutex{}

var com string
var com_lock = sync.Mutex{}

var Services chan Service

type file struct {
	*os.File
	lock sync.RWMutex
	wg   sync.WaitGroup
}

type Service struct {
	Method      byte
	Content     []byte
	Return_chan chan []byte
}

func Request(op byte, content []byte) ([]byte, error) {
	if op != Write {
		Services <- Service{Method: op, Content: nil}
	}
	return nil, nil

}

func ReadFile(name string) ([]byte, error) {
	map_lock.RLock()
	v, ok := opened_file[name]
	map_lock.RUnlock()
	if ok {
		v.wg.Add(1)
		defer v.wg.Done()
		v.lock.RLock()
		defer v.lock.RUnlock()
		return ioutil.ReadAll(v)
	}
	fs, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	wg := &sync.WaitGroup{}
	f := file{File: fs, lock: sync.RWMutex{}, wg: *wg}
	map_lock.Lock()
	opened_file[name] = &f
	map_lock.Unlock()
	wg.Add(1)
	f.lock.RLock()
	b, err := ioutil.ReadAll(fs)
	fs.Seek(0, 0)
	f.lock.RUnlock()
	go func() {
		time.Sleep(10000000000) // 10 seconds
		wg.Wait()
		map_lock.Lock()
		delete(opened_file, name)
		map_lock.Unlock()
	}()
	wg.Done()
	return b, err
}

func WriteFile(name string, data []byte) (int, error) {
	map_lock.RLock()
	v, ok := opened_file[name]
	map_lock.RUnlock()
	if ok {
		v.wg.Add(1)
		defer v.wg.Done()
		v.lock.Lock()
		defer v.lock.Unlock()
		return v.Write(data)
	}
	fs, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	wg := &sync.WaitGroup{}
	f := file{File: fs, lock: sync.RWMutex{}, wg: *wg}
	map_lock.Lock()
	opened_file[name] = &f
	map_lock.Unlock()
	wg.Add(1)
	f.lock.Lock()
	b, err := fs.Write(data)
	f.lock.Unlock()
	go func() {
		time.Sleep(10000000000) // 10 seconds
		wg.Wait()
		map_lock.Lock()
		delete(opened_file, name)
		map_lock.Unlock()
	}()

	wg.Done()

	return b, err
}

func OverWriteFile(name string, data []byte) (int, error) {
	map_lock.RLock()
	v, ok := opened_file[name]
	map_lock.RUnlock()
	if ok {
		v.wg.Add(1)
		defer v.wg.Done()
		v.lock.Lock()
		defer v.lock.Unlock()
		v.Truncate(0)
		v.Seek(0, 0)
		return v.Write(data)
	}
	fs, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	wg := &sync.WaitGroup{}
	f := file{File: fs, lock: sync.RWMutex{}, wg: *wg}
	map_lock.Lock()
	opened_file[name] = &f
	map_lock.Unlock()
	wg.Add(1)
	f.lock.Lock()
	b, err := fs.Write(data)
	f.lock.Unlock()
	go func() {
		time.Sleep(10000000000) // 10 seconds
		wg.Wait()
		map_lock.Lock()
		delete(opened_file, name)
		map_lock.Unlock()
	}()

	wg.Done()

	return b, err
}

func CreateFile(name string, data []byte) (int, error) {
	map_lock.RLock()
	v, ok := opened_file[name]
	map_lock.RUnlock()
	if ok {
		v.wg.Add(1)
		defer v.wg.Done()
		v.lock.Lock()
		defer v.lock.Unlock()
		v.Truncate(0)
		v.Seek(0, 0)
		return v.Write(data)
	}
	fs, err := os.Create(name)
	if err != nil {
		return 0, err
	}
	wg := &sync.WaitGroup{}
	f := file{File: fs, lock: sync.RWMutex{}, wg: *wg}
	map_lock.Lock()
	opened_file[name] = &f
	map_lock.Unlock()
	wg.Add(1)
	f.lock.Lock()
	b, err := fs.Write(data)
	f.lock.Unlock()
	go func() {
		time.Sleep(10000000000) // 10 seconds
		wg.Wait()
		map_lock.Lock()
		delete(opened_file, name)
		map_lock.Unlock()
	}()

	wg.Done()

	return b, err
}
