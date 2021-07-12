package main

import "net"

type data_replicates_group struct {
	connections []net.Conn
	size        int
}

func (gp *data_replicates_group) ActiveCount() int {
	count := 0
	for i, v := range gp.connections {
		_, err := v.Read(make([]byte, 0))
		if err != nil {
			gp.connections = append(gp.connections[:i], gp.connections[i+1:]...)
		} else {
			count += 1
		}
	}
	return count
}

var data_replicates_groups = map[int][]net.Conn{}

func add_replicate_group(conns []net.Conn) {
	data_replicates_groups[len(data_replicates_groups)] = conns
}

func track_connection(conn net.Conn) {

}
