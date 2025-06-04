package heapset

import (
	"fmt"
	"math"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

const (
	termW = 78 // default terminal width
)

func (hs *HeapSet) Print(op types.Output, v any) {
	ctx := buffer.Context{
		Line: true,
		Wrap: false,
		X:    0,
		Y:    0,
		W:    math.MaxInt,
		H:    math.MaxInt,
	}

	hs.RLock()

	for _, h := range hs.heaps {
		if h.Type == types.Stdin {
			continue
		}

		ctx.Heap = h.Ensure()

		switch op {
		case types.File:
			printFile(&ctx)
		case types.Grep:
			printGrep(&ctx)
		case types.Hex:
			printHex(&ctx)
		case types.Hash:
			printHash(&ctx, v.(string))
		case types.Stats:
			printStats(&ctx)
		case types.Strings:
			printStrings(&ctx, v.(int))
		}
	}

	hs.RUnlock()
}

func printFile(ctx *buffer.Context) {
	if ctx.Heap.Len() == 0 {
		return // ignore empty files
	}

	fmt.Print(string(*ctx.Heap.MMap()))
}

func printGrep(ctx *buffer.Context) {
	if ctx.Heap.Len() == 0 {
		return // ignore empty files
	}

	for l := range buffer.Text(ctx).Lines {
		fmt.Printf("%s:%s\n", ctx.Heap.String(), l)
	}
}

func printHex(ctx *buffer.Context) {
	ctx.W = termW

	fmt.Println(text.Title(ctx.Heap.String(), ctx.W))
	for l := range buffer.Hex(ctx).Lines {
		fmt.Println(l)
	}
}

func printHash(ctx *buffer.Context, sum string) {
	buf, err := ctx.Heap.HashSum(sum)

	if err != nil {
		sys.Exit(err)
	}

	fmt.Printf("%x  %s\n", buf, ctx.Heap.String())
}

func printStats(ctx *buffer.Context) {
	fmt.Printf("%8dL %8dB  %s\n",
		ctx.Heap.Count(),
		ctx.Heap.Len(),
		ctx.Heap.String(),
	)
}

func printStrings(ctx *buffer.Context, min int) {
	fmt.Println(text.Title(ctx.Heap.String(), termW))
	for s := range ctx.Heap.Strings(min) {
		fmt.Printf("%08x  %s\n", s.Off, strings.TrimSpace(s.Str))
	}
}
