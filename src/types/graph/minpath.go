package graph

import (
	"errors"
)

type Path struct {
	Path []Edge
	Cost int64
}

type joinPath struct {
	next  *Vertex
	paths []*Path
	cost  int64
}

func ShortestPath(from, to *Vertex) (Path, error) {
	from.Lock.RLock()
	a := *from
	from.Lock.RUnlock()

	costs := map[*Vertex]*Path{} // from * towards key weight
	costs[from] = &Path{
		Path: []Edge{},
		Cost: 0,
	}
	items_store := []joinPath{joinPath{next: a}}

	endedPaths := []joinPath{}

	var lowest_score int64 = 9223372036854775807

	for {
		if len(items_store) == 0 {
			break
		}
		items := []joinPath{}
		for _, v := range items_store {
			for _, e := range v.next.Edges {
				if e.To == from {
					continue
				}

				new := *costs[e.From]
				new.Cost += e.Weight
				new.Path = append(new.Path, v)

				if p, ok := costs[e.To]; ok {
					if p.Cost > new.Cost {
						*(costs[e.To]) = new
					}
				} else {
					costs[e.To] = &new

					var c int64 = 0
					for _, i := range v.paths {
						c += i.Cost
					}

					p := joinPath{
						next:  e.To,
						paths: append(v.paths, &new),
						cost:  c + new.Cost,
					}
					if e.To == to {
						if p.cost < lowest_score {
							lowest_score = p.cost
						}
						endedPaths = append(endedPaths, p)
					} else if !(p.cost > lowest_score) {
						items = append(items, p) // items extend the loop
					}

				}
			}
		}
		items_store = items
	}
	if len(endedPaths) == 0 {
		return Path{}, errors.New("no paths found")
	}
	var result = Path{
		Path: []Edge{},
		Cost: 0,
	}

	for _, v := range endedPaths {
		var c int64 = 0
		var E = []Edge{}
		for _, v := range v.paths {
			E = append(E, v.Path)
			c += v.Cost
		}
		if result.Cost > c {
			result.Path = E
			result.Cost = c
		}
	}
	return result, nil
}
