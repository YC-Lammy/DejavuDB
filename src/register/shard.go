package register

import (
	"sync"
)

type shard struct {
	size      int
	conn      net.Conn
	mem_load  int
	cpu_load  int
	disk_load int
	mem_size  int
	disk_size int
}

var Shard = map[uint16]*shard{}
var shardlock = sync.Mutax{}

func Get(id uint16)(*shard, error){
	Lock.Lock()
	if v,ok:=Shard[id];ok{
		lock.Unlock()
		return v, nil
	}
	lock.Unlock()
	return nil, errors.New("key id not exist")
}

func Set(id uint16, data *shard)error{
	lock.Lock()
	if data != nil{
		Shard[id] = data
		lock.Unlock()
		return nil
	}
	return errors.New("data cannot be nil")
}