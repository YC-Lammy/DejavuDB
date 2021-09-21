package graph

import (
	"fmt"
	"sync"

	"github.com/golang/snappy"
)

type Edge struct {
	Id     uint64
	Label  string
	Weight int16
	From   *Vertex
	To     *Vertex
	Next   *Edge
}

// Node a single node that composes the tree
type Vertex struct {
	Id     uint64
	Values map[string]string //key value store
	Edges  []*Edge
	Lock   sync.RWMutex
}

func NewVertex(Values map[string]string) *Vertex {
	return &Vertex{
		Lock: sync.RWMutex{},
	}
}

func (n *Vertex) AddEdge(e *Edge) {
	if n.Edges == nil {
		n.Edges = []*Edge{}
	}
	n.Lock.Lock()
	n.Edges = append(n.Edges, e)
	n.Lock.Unlock()
}

func (n *Vertex) AddField(key, data string) {
	n.Values[key] = string(snappy.Encode(nil, []byte(data)))
}
func (n *Vertex) GetField(key string) (string, error) {
	a, err := snappy.Decode(nil, []byte(n.Values[key]))
	if err != nil {
		return "", err
	}
	return string(a), nil
}
func (n *Vertex) String() string {
	return fmt.Sprintf("%v", n.Values)
}

type Graph struct {
	Nodes             map[uint64]*Vertex
	Lock              sync.RWMutex
	Edge_count        uint64
	Edge_count_Lock   sync.RWMutex
	Vertex_count      uint64
	Vertex_count_Lock sync.RWMutex
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(n *Vertex) {
	g.Edge_count_Lock.Lock()
	g.Vertex_count += 1
	n.Id = g.Vertex_count
	g.Edge_count_Lock.Unlock()

	g.Lock.Lock()
	g.Nodes[g.Vertex_count] = n
	g.Lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(n1, n2 *Vertex, label string, weight int16) {
	if n1.Edges == nil {
		n1.Edges = []*Edge{}
	}
	if n2.Edges == nil {
		n2.Edges = []*Edge{}
	}
	g.Edge_count_Lock.Lock()
	g.Edge_count += 1
	e := g.Edge_count
	g.Edge_count_Lock.Unlock()

	a := Edge{
		Id:     e,
		Label:  label,
		Weight: weight,
		From:   n1,
		To:     n2,
	}
	n1.AddEdge(&a)
	n2.AddEdge(&a)
}

func (g *Graph) String() string {
	g.Lock.RLock()
	s := ""
	for _, k := range g.Nodes {
		s += k.String() + " -> "
		near := k.Edges
		for j := 0; j < len(near); j++ {
			if near[j].To == k {
				continue
			}
			s += near[j].To.String() + " "
		}
		s += "\n"
	}
	g.Lock.RUnlock()
	return s
}
