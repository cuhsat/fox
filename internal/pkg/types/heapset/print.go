package heapset

import (
	"fmt"
	"math"

	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/text"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/buffer"
)

func (hs *HeapSet) Print(op types.Output, sum string) {
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

		ctx.Heap = h.Ensure().Filter()

		switch op {
		case types.File:
			printFile(&ctx)
		case types.Grep:
			printGrep(&ctx)
		case types.Hex:
			printHex(&ctx)
		case types.Hash:
			printHash(&ctx, sum)
		case types.Count:
			printCount(&ctx)
		}
	}

	hs.RUnlock()
}

func printFile(ctx *buffer.Context) {
	if ctx.Heap.Lines() != 0 {
		fmt.Print(buffer.Text(ctx))
	}
}

func printGrep(ctx *buffer.Context) {
	if ctx.Heap.Lines() != 0 {
		for _, tl := range buffer.Text(ctx).Lines {
			fmt.Printf("%s:%s\n", ctx.Heap.String(), tl)
		}
	}
}

func printHex(ctx *buffer.Context) {
	ctx.W = 78 // use default terminal width

	fmt.Println(text.Title(ctx.Heap.String(), ctx.W))
	fmt.Println(buffer.Hex(ctx))
}

func printHash(ctx *buffer.Context, sum string) {
	buf, err := ctx.Heap.Hashsum(sum)

	if err != nil {
		sys.Exit(err)
	}

	fmt.Printf("%x  %s\n", buf, ctx.Heap.String())
}

func printCount(ctx *buffer.Context) {
	fmt.Printf("%8dL %8dB  %s\n", ctx.Heap.Size(), ctx.Heap.Size(), ctx.Heap.String())
}
