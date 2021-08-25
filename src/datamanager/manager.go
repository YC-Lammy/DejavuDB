package datamanager

import (
	"net"
)

const (
	Follower  = 0x00
	Candidate = 0x01
	Leader    = 0x02
)

type shard struct {
	Conn *net.Conn
	Id   int
}

type location struct {

	//dataGroup
}

var DataMap = map[string]*location{}
