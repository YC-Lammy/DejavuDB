package graph

import (
	"fmt"
	"sync"
)

type Edge struct {
	Id     uint64
	Label  [16]byte
	Weight int16
	From   *Node
	To     *Node
	Next   *Edge //linked list
}

// Node a single node that composes the tree
type Node struct {
	Id     uint64
	Name   string
	Values map[string]string //key value store
	Edges  []*Node
	Lock   sync.Mutex
}

func (n *Node) AddEdge(e *Node) {
	if n.Edges == nil {
		n.Edges = []*Node{}
	}
	n.Lock.Lock()
	n.Edges = append(n.Edges, e)
	n.Lock.Unlock()
}

func (n *Node) AddField(key, data string) {
	n.Values[key] = data
}
func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Values)
}

type Graph struct {
	Nodes map[string]*Node
	Lock  sync.RWMutex
}

func (g *Graph) AddNode(n *Node) {
	g.Lock.Lock()
	g.Nodes[n.Name] = n
	g.Lock.Unlock()
}

// ItemGraph the Items graph
type ItemGraph struct {
	nodes []*Node
	lock  sync.RWMutex
}

// AddNode adds a node to the graph
func (g *ItemGraph) AddNode(n *Node) {
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) AddEdge(n1, n2 *Node) {
	if n1.Edges == nil {
		n1.Edges = []*Node{}
	}
	if n2.Edges == nil {
		n2.Edges = []*Node{}
	}
	n1.Lock.Lock()
	n1.Edges = append(n1.Edges, n2)
	n1.Lock.Unlock()
	n2.Lock.Lock()
	n2.Edges = append(n2.Edges, n1)
	n2.Lock.Unlock()
}

func (g *ItemGraph) String() string {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.nodes[i].Edges
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	g.lock.RUnlock()
	return s
}
