package graph

import "container/heap"

type a_star_item struct {
	now               *Vertex
	prev              *Vertex
	distanceFromStart int64
	priority          int32
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

func A_Star(from, to *Vertex, heuristic func(key, endKey string) int) ([]string, error) {

	// priorityQueue for nodes that have not yet been visited (open nodes)
	openQueue := &a_star_priority_queue{}
	openList := map[uint64]*a_star_item{}
	closedList := map[uint64]*a_star_item{}

	item := &a_star_item{from, nil, 0, 0, 0}
	openList[from.Id] = item

	heap.Push(openQueue, item)

	for len(*openQueue) >0 {
		current := heap.Pop(openQueue).()
	}
}
