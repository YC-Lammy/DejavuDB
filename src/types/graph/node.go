package graph

func (n *Node) AddEdge(e *Node) {
	if n.Edges == nil {
		n.Edges = []*Node{}
	}
	n.Lock.Lock()
	n.Edges = append(n.Edges, e)
	n.Lock.Unlock()
}
