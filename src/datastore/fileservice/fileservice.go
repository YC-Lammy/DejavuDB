package fileservice

import (
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var opened_file map[string]*file
var map_lock = &sync.RWMutex{}

var com string
var com_lock = &sync.Mutex{}

type file struct {
	*os.File
	lock *sync.RWMutex
	wg   *sync.WaitGroup
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
	l := sync.RWMutex{}
	wg := &sync.WaitGroup{}
	f := file{File: fs, lock: &l, wg: wg}
	map_lock.Lock()
	opened_file[name] = &f
	map_lock.Unlock()
	wg.Add(1)
	l.RLock()
	b, err := ioutil.ReadAll(fs)
	l.RUnlock()
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
	l := sync.RWMutex{}
	wg := &sync.WaitGroup{}
	f := file{File: fs, lock: &l, wg: wg}
	map_lock.Lock()
	opened_file[name] = &f
	map_lock.Unlock()
	wg.Add(1)
	l.Lock()
	b, err := fs.Write(data)
	l.Unlock()
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
