package buffer

import (
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

const (
	Size = 1024
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
