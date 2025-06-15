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
	Head bool

	Nr int

	X int
	Y int
	W int
	H int
}

func (ctx *Context) Hash(suffix string) string {
	return fmt.Sprintf("%s#%d:%t:%t:%t@%d:%d|%s",
		ctx.Heap.LastFilter().Pattern,
		ctx.Heap.Len(),
		ctx.Wrap,
		ctx.Line,
		ctx.Head,
		ctx.W,
		ctx.H,
		suffix,
	)
}
