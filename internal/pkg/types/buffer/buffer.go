package buffer

import (
	"fmt"
	"math"

	"github.com/hiforensics/fox/internal/pkg/types/heap"
)

const (
	Size = 1024
)

const (
	TermW = 78
	TermH = 24
)

type Buffer struct {
	W int
	H int
}

type Line struct {
	Nr  string
	Grp int
	Str string
}

type Part struct {
	X   int
	Y   int
	Grp int
	Str string
}

type Context struct {
	Heap *heap.Heap

	Context bool
	Numbers bool
	Wrap    bool

	Nr int

	X int
	Y int
	W int
	H int
}

func NewContext(h *heap.Heap) *Context {
	return &Context{
		Context: true,
		Numbers: true,
		Wrap:    false,
		Heap:    h,
		X:       0,
		Y:       0,
		W:       math.MaxInt,
		H:       math.MaxInt,
	}
}

func (ctx *Context) Hash() string {
	return fmt.Sprintf("%s[%d:%d]#%d:%t:%t:%t@%d:%d",
		ctx.Heap.LastFilter().Pattern,
		ctx.Heap.LastFilter().Context.B,
		ctx.Heap.LastFilter().Context.A,
		ctx.Heap.Len(),
		ctx.Context,
		ctx.Numbers,
		ctx.Wrap,
		ctx.W,
		ctx.H,
	)
}
