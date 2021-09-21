package graph

import "sync"

// NodeQueue the queue of Nodes
type NodeQueue struct {
	items []Vertex
	lock  sync.RWMutex
}

// New creates a new NodeQueue
func (s *NodeQueue) New() *NodeQueue {
	s.lock.Lock()
	s.items = []Vertex{}
	s.lock.Unlock()
	return s
}

// Enqueue adds an Node to the end of the queue
func (s *NodeQueue) Enqueue(t Vertex) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Node from the start of the queue
func (s *NodeQueue) Dequeue() *Vertex {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:]
	s.lock.Unlock()
	return &item
}

// Front returns the item next in the queue, without removing it
func (s *NodeQueue) Front() *Vertex {
	s.lock.RLock()
	item := s.items[0]
	s.lock.RUnlock()
	return &item
}

// IsEmpty returns true if the queue is empty
func (s *NodeQueue) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items) == 0
}

// Size returns the number of Nodes in the queue
func (s *NodeQueue) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

func (g *Graph) Traverse(f func(*Vertex)) {
	g.Lock.RLock()
	q := NodeQueue{}
	q.New()
	n := g.Nodes[0]
	q.Enqueue(*n)
	visited := make(map[*Vertex]bool)
	for {
		if q.IsEmpty() {
			break
		}
		node := q.Dequeue()
		visited[node] = true
		near := node.Edges

		for i := 0; i < len(near); i++ {
			j := near[i]
			if j.To == node {
				continue
			}
			if !visited[j.To] {
				q.Enqueue(*j.To)
				visited[j.To] = true
			}
		}
		if f != nil {
			f(node)
		}
	}
	g.Lock.RUnlock()
}
