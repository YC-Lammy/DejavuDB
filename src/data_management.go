package main

import "net"

type data_replicates_group struct {
	connections []net.Conn
	size        int
	addrs       []string
	load        int // the highest percentage in the group
}

var data_replicates_groups = map[int]*data_replicates_group{}

var data_group_addrs = map[int][]string{} // group number, addrs

func (gp *data_replicates_group) ActiveCount() int {
	count := 0
	for i, v := range gp.connections {
		_, err := v.Read(make([]byte, 0)) // remove closed connections
		if err != nil {
			gp.connections = append(gp.connections[:i], gp.connections[i+1:]...)
		} else {
			count += 1
		}
	}
	return count
}

func (gp *data_replicates_group) Update() {
	addrs := []string{}
	for _, v := range gp.connections {
		addrs = append(addrs, v.RemoteAddr().String())
	}
	gp.addrs = addrs
}

func add_replicate_group(conns []net.Conn) {
	addrs := []string{}
	for _, v := range conns {
		addrs = append(addrs, v.RemoteAddr().String())
	}
	data_replicates_groups[len(data_replicates_groups)] = &data_replicates_group{connections: conns, addrs: addrs}
}

func remove_replicate_group(groupnumber int) {
	delete(data_replicates_groups, groupnumber)
}

func track_connection(conn net.Conn) {

}
