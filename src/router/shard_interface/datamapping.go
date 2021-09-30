package shard_interface

import (
	"src/config"
	"src/network"
	"strings"
	"sync"
)

var DataMap = map[string]*ReplicaGroup{}
var DataMapLock = sync.RWMutex{}

var ReplicaGroups = map[uint64]*ReplicaGroup{}

var groupCount uint64

type ReplicaGroup struct {
	Id        uint64
	Shards    []*Shard_conn
	Ram_load  uint8
	Disk_load uint8
}

func NewReplicateGroup(shards ...*Shard_conn) *ReplicaGroup {
	groupCount += 1
	a := groupCount
	r := &ReplicaGroup{
		Id:     a,
		Shards: shards,
	}
	ReplicaGroups[a] = r
	return r
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

func (r *ReplicaGroup) Send(ProcessId uint64, Script []byte) (err error) {
	for _, v := range r.Shards {
		_, e := network.Router_Shard_Send(v.Conn, uint64(config.ID), ProcessId, Script)
		if e != nil {
			err = e
		}
	}
	return
}
