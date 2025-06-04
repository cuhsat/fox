package buffer

import (
	"fmt"

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

	Nr int

	X int
	Y int
	W int
	H int
}

func (ctx *Context) Hash() string {
	return fmt.Sprintf("%s#%d-%t-%t@%d:%d",
		ctx.Heap.LastFilter().Pattern,
		ctx.Heap.Len(),
		ctx.Wrap,
		ctx.Line,
		ctx.W,
		ctx.H,
	)
}
