package network

import (
	"errors"
	"net"
	"sync"
)

var Shard_connected uint64 = 0
var Router_connected uint64 = 0
var Client_connected uint64 = 0

type shard struct {
	Id        uint32
	Size      int
	Conn      net.Conn
	Mem_load  int
	Cpu_load  int
	Disk_load int
	Mem_size  int
	Disk_size int
}

var Shards = map[uint16]*shard{}
var shardlock = sync.Mutex{}

func GetShard(id uint16) (*shard, error) {
	shardlock.Lock()
	if v, ok := Shards[id]; ok {
		shardlock.Unlock()
		return v, nil
	}
	shardlock.Unlock()
	return nil, errors.New("key id not exist")
}

func SetShard(id uint16, data *shard) error {
	if data != nil {
		shardlock.Lock()
		Shards[id] = data
		shardlock.Unlock()
		return nil
	}
	return errors.New("data cannot be nil")
}
