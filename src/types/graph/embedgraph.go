package graph

import (
	"crypto/rand"
	"math"
	"math/big"
)

type EmbedGraph struct {
	Sections [][]*EmbedSection
	X        uint32 // default 1000000
	Y        uint32 // default 1000000
}

type EmbedSection struct {
	Vertexs []*EmbedVertex
}

type EmbedVertex struct {
	Vertex
	X             uint32
	Y             uint32
	Applied_force uint64 // when more force is applied, more friction
}

func NewEmbedVertex() *EmbedVertex {

	x, _ := rand.Int(rand.Reader, big.NewInt(4294967295))
	y, _ := rand.Int(rand.Reader, big.NewInt(4294967295))
	return &EmbedVertex{Vertex: *NewVertex(), X: uint32(x.Uint64()), Y: uint32(y.Uint64())}
}

func (e *EmbedVertex) ApplyForce(weight uint32, fromX, fromY uint32) {

	var p float64
	if e.Applied_force < uint64(weight) {
		p = 1
	} else {
		p = float64(weight) / float64(e.Applied_force)
	}

	force := p * float64(weight)

	var lenX uint32
	var lenY uint32

	if fromX > e.X {
		lenX = fromX - e.X
	} else {
		lenX = e.X - fromX
	}

	if fromY > e.Y {
		lenY = fromY - e.Y
	} else {
		lenY = e.Y - fromY
	}

	hy := math.Sqrt(float64((lenX * lenX) + (lenY * lenY)))

	e.Applied_force += weight

}
