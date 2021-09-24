package network

import (
	"net"

	"src/config"
	"src/meta"

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

func Router_Shard_Send(conn net.Conn, Process_Id uint64, script []byte) (int, error) {
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
		Role:       role,
		Id:         meta.Id,
		Process_Id: Process_Id,
		Script:     script,
	})
	if err != nil {
		return 0, err
	}
	return Send(conn, b)
}

func Router_Shard_Recv(conn net.Conn) (Package, error) {
	var p = Package{}
	b, err := Receive(conn)
	if err != nil {
		return p, err
	}
	err = cbor.Unmarshal(b, &p)
	return p, err
}
