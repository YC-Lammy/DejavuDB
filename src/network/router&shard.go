package network

import (
	"../config"
	"../meta"
	cbor "github.com/fxamacker/cbor/v2"
)

const (
	router_role     = 0x00
	shard_role      = 0x01
	standalone_role = 0x02
)

type Package struct {
	Role       byte
	Id         uint64
	Process_Id uint64
	Script     []byte
}

func Router_Shard_Send(Id, Process_Id uint64, script []byte) {
	var role byte
	switch config.Role {
	case "router":
		role = router_role
	case "shard":
		role = shard_role
	case "standalone":
		role = standalone_role
	}
	b, err := cbor.Marshal(Package{
		Role: role,
		Id:   meta.Id,
	})
}
