package graph

import (
	"errors"
	"strings"
)

type Path struct {
	Path []*Edge
	Cost int64
}

func ShortestPath(from, to *Vertex, method string) (Path, error) {
	switch strings.ToLower(method) {
	case "dijkstra":
		return Dijkstra(from, to)
	}
	return Path{}, errors.New("unknown method " + method)
}

func Dijkstra(from, to *Vertex) (Path, error) {
	from.Lock.RLock()
	from.Lock.RUnlock()

	costs := map[uint64]*Path{} // from * towards key weight
	costs[from.Id] = &Path{
		Path: []*Edge{},
		Cost: 0,
	}
	items_store := []*Vertex{from}

	//endedPaths := []joinPath{}

	var lowest_score int64 = 9223372036854775807

	for {
		if len(items_store) == 0 {
			break
		}
		items := []*Vertex{}
		for _, v := range items_store {
			for _, e := range v.Edges {
				if v == e.To {
					continue
				}

				new := Path{}
				if n, ok := costs[e.From.Id]; ok {
					new = *n
					new.Cost = new.Cost + int64(e.Weight)
					new.Path = append(new.Path, e)

				} else {
					return Path{}, errors.New("error processing node")
				}

				if p, ok := costs[e.To.Id]; ok {
					if p.Cost > new.Cost {

						if e.To.Id == to.Id {

							if new.Cost < lowest_score {
								lowest_score = new.Cost
								costs[e.To.Id] = &new
							}
							//endedPaths = append(endedPaths, p)
						} else if new.Cost < lowest_score {
							costs[e.To.Id] = &new
							items = append(items, e.To) // items extend the loop
						}
					}
				} else {
					/*
						var c int64 = 0
						for _, i := range v.paths {
							c += i.Cost
						}
					*/

					if e.To.Id == to.Id {

						if new.Cost < lowest_score {
							lowest_score = new.Cost
							costs[e.To.Id] = &new
						}
						//endedPaths = append(endedPaths, p)
					} else if new.Cost < lowest_score {
						costs[e.To.Id] = &new
						items = append(items, e.To) // items extend the loop
					}

				}
			}
		}
		items_store = items
	}
	if v, ok := costs[to.Id]; ok {
		return *v, nil
	} else {
		return Path{}, errors.New("no path found")
	}
}
