package graph

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_shortest_bfs(t *testing.T) {
	var num_node = 4096
	g := NewGraph()
	a := NewVertex(map[string]string{})
	g.AddNode(a)
	var last *Vertex = a
	for i := 0; i < num_node; i++ {
		b := NewVertex(map[string]string{})
		g.AddNode(b)
		rand.Seed(time.Now().Unix())
		g.AddEdge(last, b, "", 1)
		last = b
	}
	g.AddEdge(a, g.Nodes[293], "", 290)
	T := time.Now()
	p, err := ShortestPath(a, last, "dijkstra")
	fmt.Println(time.Now().Sub(T).String())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p.Cost)
	//g.Traverse(func(v *Vertex) {})
	//fmt.Println(g.String())
}
