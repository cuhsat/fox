package buffer

import (
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/smap"
)

const (
	Size = 1024
)

var (
	Cache = make(map[string]*smap.SMap, 256)
)

type Buffer struct {
	W int
	H int
}

type Line struct {
	Nr  string
	Str string
}

type Part struct {
	X   int
	Y   int
	Str string
}

type Context struct {
	Heap *heap.Heap

	Wrap bool
	Line bool

	X int
	Y int
	W int
	H int
}
