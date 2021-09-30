package shard_interface

type ReplicaGroup struct {
	Id     uint64
	Shards []*Shard_conn
}
