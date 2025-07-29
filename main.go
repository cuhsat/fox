package main

import (
	"os"
	"runtime/debug"

	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/arg"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
	"github.com/hiforensics/fox/internal/pkg/user/bag"
	"github.com/hiforensics/fox/internal/pkg/user/history"
)

func main() {
	args := arg.GetArgs()

	sys.Setup()

	defer func() {
		if err := recover(); err != nil {
			sys.DumpErr(err, debug.Stack())
			sys.Print(err)
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
