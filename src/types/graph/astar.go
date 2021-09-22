package graph

import (
	"container/heap"
	"errors"
)

type a_star_item struct {
	edge              *Edge // from is previous, to is now
	prev              *Edge
	distanceFromStart int64
	priority          int64
	index             int64
}
type a_star_priority_queue []*a_star_item

func (a a_star_priority_queue) Len() int {
	return len(a)
}
func (a a_star_priority_queue) Less(i, j int) bool {
	return a[i].priority < a[j].priority
}
func (a a_star_priority_queue) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
	a[i].index = int64(i)
	a[j].index = int64(j)
}
func (pq *a_star_priority_queue) Push(item interface{}) {
	a := item.(*a_star_item)
	a.index = int64(len(*pq))
	*pq = append(*pq, a)
}
func (pq *a_star_priority_queue) Pop() interface{} {
	item := (*pq)[len(*pq)-1]
	*pq = (*pq)[0 : len(*pq)-1]
	return item
}

func A_Star(from, to *Vertex, heuristic func(key, endKey uint64) int64) (path []*Edge, err error) {

	// priorityQueue for nodes that have not yet been visited (open nodes)
	openQueue := &a_star_priority_queue{}
	openList := map[uint64]*a_star_item{} // key is vertex id
	closedList := map[uint64]*a_star_item{}

	item := &a_star_item{
		&Edge{From: nil, To: from},
		nil,
		0, 0, 0}
	openList[from.Id] = item

	heap.Push(openQueue, item)

	for len(*openQueue) > 0 {
		current_edge := heap.Pop(openQueue).(*a_star_item).edge
		current := current_edge.To
		// current node was now visited; add to closed list
		closedList[current.Id] = openList[current.Id]
		delete(openList, current.Id)

		// end node found?
		if current == to {
			// build path
			for current_edge != nil {
				path = append(path, current_edge)
				a := closedList[current.Id]
				current = a.edge.From
				current_edge = a.prev
			}

			return path, nil
		}
		// saved here for easy usage in following loop
		distance := closedList[current.Id].distanceFromStart

		for _, edge := range current.Edges {
			neighbor := edge.To
			if _, ok := closedList[neighbor.Id]; ok {
				continue
			}

			distanceToNeighbor := distance + int64(edge.Weight)

			// skip neighbors that already have a better path leading to them
			if md, ok := openList[neighbor.Id]; ok {
				if md.distanceFromStart < distanceToNeighbor {
					continue
				} else {
					heap.Remove(openQueue, int(md.index))
				}
			}

			item := &a_star_item{
				edge,
				current_edge,
				int64(distanceToNeighbor),
				int64(distanceToNeighbor + heuristic(neighbor.Id, to.Id)), // estimate (= priority)
				0,
			}

			// add neighbor node to list of open nodes
			openList[neighbor.Id] = item

			// push into priority queue
			heap.Push(openQueue, item)
		}
	}
	return nil, errors.New("no path found")
}
