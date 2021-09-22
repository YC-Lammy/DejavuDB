package graph

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_shortest_bfs(t *testing.T) {
	var num_node = 40960
	g := NewGraph()
	a := NewVertex()
	g.AddNode(a)
	var last *Vertex = a
	for i := 0; i < num_node; i++ {
		b := NewVertex()
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

func Test_A_star(t *testing.T) {
	var num_node = 40960
	g := NewGraph()
	a := NewVertex()
	g.AddNode(a)
	var last *Vertex = a
	for i := 0; i < num_node; i++ {
		b := NewVertex()
		g.AddNode(b)
		rand.Seed(time.Now().Unix())
		g.AddEdge(last, b, "", 1)
		last = b
	}
	g.AddEdge(a, g.Nodes[293], "", 290)
	T := time.Now()
	p, err := A_Star(a, last, func(a, b uint64) int64 { return int64(a - b) })
	fmt.Println(time.Now().Sub(T).String())
	if err != nil {
		fmt.Println(err)
		return
	}
	var f int64
	for _, v := range p {
		f += int64(v.Weight)
	}
	fmt.Println(f)
}
