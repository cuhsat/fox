package heapset

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/hiforensics/fox/internal/pkg/arg"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/text"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/buffer"
)

const (
	termW = 78 // default terminal width
)

func (hs *HeapSet) Deflate(args *arg.Args) {
	hs.RLock()

	for _, h := range hs.heaps {
		root := args.Run.Deflate

		if root == "-" {
			name := filepath.Base(h.Base)
			root = name[0 : len(name)-len(filepath.Ext(name))]
		}

		// convert to relative path
		path := h.Title

		if h.Type == types.Deflate {
			path = path[len(h.Base)+1:]
		} else {
			path = filepath.Base(path)
		}

		// create (sub)folders
		if sub := filepath.Dir(path); len(sub) > 0 {
			sub = filepath.Join(root, sub)

			err := os.MkdirAll(sub, 0700)

			if err != nil {
				sys.Exit(err)
			}
		}

		path = filepath.Join(root, path)

		if !args.Print.NoFile {
			fmt.Printf("Deflate %s\n", path)
		}

		err := os.WriteFile(path, *h.Ensure().MMap(), 0600)

		if err != nil {
			sys.Exit(err)
		}
	}

	hs.RUnlock()

	if l := sys.Log.Consume(); len(l) > 0 {
		fmt.Fprintln(os.Stderr, l)
	}

	fmt.Printf("%d file(s) written\n", hs.Len())
}

func (hs *HeapSet) Print(args *arg.Args) {
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

		switch args.Run.Mode {
		case types.File:
			printFile(&ctx)
		case types.Grep:
			printGrep(&ctx)
		case types.Hex:
			printHex(&ctx)
		case types.Hash:
			printHash(&ctx, args.Run.Value.(string))
		case types.Stats:
			printStats(&ctx)
		case types.Strings:
			printStrings(&ctx, args.Run.Value.(int))
		}
	}

	hs.RUnlock()

	if l := sys.Log.Consume(); len(l) > 0 {
		fmt.Fprintln(os.Stderr, l)
	}
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
