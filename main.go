package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/cuhsat/fox/internal/fox/ui"
	"github.com/cuhsat/fox/internal/pkg/arg"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
	"github.com/cuhsat/fox/internal/pkg/user/bag"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

func main() {
	args := arg.GetArgs()

	sys.Setup()

	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			sys.DumpErr(err, debug.Stack())
		}

		sys.Log.Close()
	}()

	_ = os.Remove(sys.Dump)

	hs := heapset.New(args.Args)
	defer hs.ThrowAway()

	if args.Print.Active {
		hs.Print(args.Print)
	} else if len(args.Deflate) > 0 {
		hs.Deflate(args.Deflate)
	} else {
		hi := history.New()
		defer hi.Close()

		bg := bag.New(args.Bag)
		defer bg.Close()

		fx := ui.New(args.UI)
		defer fx.Close()

		fx.Run(hs, hi, bg)
	}
}
