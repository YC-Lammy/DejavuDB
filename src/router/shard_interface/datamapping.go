package shard_interface

import (
	"src/network"
	"strings"
	"sync"
)

var DataMap = map[string]*ReplicaGroup{}
var DataMapLock = sync.RWMutex{}

type ReplicaGroup struct {
	Id        uint64
	Shards    []*Shard_conn
	Ram_load  uint8
	Disk_load uint8
}

func GroupByKey(key string) *ReplicaGroup {
	keys := strings.Split(key, ".")
	k := ""
	for _, v := range keys {
		k += "." + v
		DataMapLock.RLock()
		v, ok := DataMap[k]
		DataMapLock.RUnlock()
		if ok {
			return v
		}
	}
	return nil
}

func (r *ReplicaGroup) Send() {
	network.Router_Shard_Send()
}
