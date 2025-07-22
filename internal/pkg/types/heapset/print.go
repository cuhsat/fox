package heapset

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/arg"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
)

const (
	termW = 78 // default terminal width
)

func (hs *HeapSet) Deflate(path string) {
	args := arg.GetArgs()

	hs.RLock()

	for _, h := range hs.heaps {
		b := filepath.Base(h.Base)

		dir := b[0 : len(b)-len(filepath.Ext(b))]

		// convert to relative path
		p := h.Title

		if h.Type == types.Deflate {
			p = p[len(h.Base)+1:]
		} else {
			p = filepath.Base(p)
		}

		// create (sub)folders
		if d := filepath.Dir(p); len(d) > 0 {
			d = filepath.Join(path, dir, d)

			err := os.MkdirAll(d, 0700)

			if err != nil {
				sys.Exit(err)
			}
		}

		p = filepath.Join(path, dir, p)

		if !args.Print.NoFile {
			fmt.Printf("Deflate %s\n", p)
		}

		err := os.WriteFile(p, *h.Ensure().MMap(), 0600)

		if err != nil {
			sys.Exit(err)
		}
	}

	hs.RUnlock()

	fmt.Printf("%d file(s) written\n", hs.Len())
}

func (hs *HeapSet) Print(args arg.ArgsPrint) {
	ctx := buffer.Context{
		Context: true,
		Numbers: true,
		Wrap:    false,
		X:       0,
		Y:       0,
		W:       math.MaxInt,
		H:       math.MaxInt,
	}

	hs.RLock()

	for _, h := range hs.heaps {
		if h.Type == types.Stdin {
			continue
		}

		ctx.Heap = h.Ensure()

		switch args.Mode {
		case types.File:
			printFile(&ctx)
		case types.Grep:
			printGrep(&ctx)
		case types.Hex:
			printHex(&ctx)
		case types.Hash:
			printHash(&ctx, args.Value.(string))
		case types.Stats:
			printStats(&ctx)
		case types.Strings:
			printStrings(&ctx, args.Value.(int))
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

	a := arg.GetArgs().Print

	if !a.NoFile {
		fmt.Println(text.Title(ctx.Heap.String(), termW))
	}

	for l := range buffer.Text(ctx).Lines {
		if l.Nr == "--" {
			if !a.NoLine {
				fmt.Println("--")
			}
		} else {
			if !a.NoLine {
				fmt.Printf("%s:%s\n", l.Nr, l)
			} else {
				fmt.Println(l)
			}
		}
	}
}

func printHex(ctx *buffer.Context) {
	ctx.W = termW

	if !arg.GetArgs().Print.NoFile {
		fmt.Println(text.Title(ctx.Heap.String(), ctx.W))
	}

	for l := range buffer.Hex(ctx).Lines {
		fmt.Println(l)
	}
}

func printHash(ctx *buffer.Context, algo string) {
	sum, err := ctx.Heap.HashSum(algo)

	if err != nil {
		sys.Exit(err)
	}

	switch strings.ToLower(algo) {
	case "sdhash": // string results
		fmt.Printf("%s  %s\n", sum, ctx.Heap.String())
	default:
		fmt.Printf("%x  %s\n", sum, ctx.Heap.String())
	}
}

func printStats(ctx *buffer.Context) {
	fmt.Printf("%8dL %8dB  %s\n",
		ctx.Heap.Count(),
		ctx.Heap.Len(),
		ctx.Heap.String(),
	)
}

func printStrings(ctx *buffer.Context, n int) {
	a := arg.GetArgs().Print

	if !a.NoFile {
		fmt.Println(text.Title(ctx.Heap.String(), termW))
	}

	for s := range ctx.Heap.Strings(n) {
		if !a.NoLine {
			fmt.Printf("%08x  %s\n", s.Off, strings.TrimSpace(s.Str))
		} else {
			fmt.Println(strings.TrimSpace(s.Str))
		}
	}
}
