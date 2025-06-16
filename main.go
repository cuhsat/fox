package main

import (
	"fmt"
	"os"
	"runtime/debug"

	flag "github.com/spf13/pflag"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/fox/ui"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
	"github.com/cuhsat/fox/internal/pkg/types/mode"
	"github.com/cuhsat/fox/internal/pkg/user/bag"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

type args struct {
	args []string

	// console
	print bool

	// evidence bag
	file string
	mode string
	key  string

	render struct {
		mode mode.Mode
	}

	output struct {
		mode  types.Output
		value any
	}
}

func argsParse() (a args) {
	a.render.mode = mode.Default
	a.output.mode = types.File

	// console
	flag.BoolVarP(&a.print, "print", "p", false, "print to console (no UI)")

	x := flag.BoolP("hex", "x", false, "output file in hex / start in HEX mode")
	w := flag.BoolP("counts", "w", false, "output file line and byte counts")
	s := flag.IntP("strings", "s", 0, "output file ASCII and Unicode strings")
	H := flag.StringP("hash", "H", "", "output hash sum of file")

	if *s == 0 {
		flag.Lookup("strings").NoOptDefVal = "3" // default
	}

	if len(*H) == 0 {
		flag.Lookup("hash").NoOptDefVal = heap.Sha256 // default
	}

	// file limits
	limits := types.GetLimits()

	head := flag.BoolP("head", "h", false, "limit head of file by ...")
	tail := flag.BoolP("tail", "t", false, "limit tail of file by ...")

	counts := new(types.Counts)

	flag.IntVarP(&counts.Lines, "lines", "n", 0, "number of lines")
	flag.IntVarP(&counts.Bytes, "bytes", "c", 0, "number of bytes")

	if counts.Lines == 0 {
		flag.Lookup("lines").NoOptDefVal = "10" // default
	}

	if counts.Bytes == 0 {
		flag.Lookup("bytes").NoOptDefVal = "16" // default
	}

	// line filter
	filters := types.GetFilters()

	flag.VarP(filters, "regexp", "e", "filter for lines that matches pattern")

	// evidence bag
	flag.StringVarP(&a.file, "file", "f", bag.Filename, "file name of evidence bag")
	flag.StringVarP(&a.mode, "mode", "m", "", "output mode")
	flag.StringVarP(&a.key, "key", "k", "", "key phrase for signing with HMAC")

	if len(a.mode) == 0 {
		flag.Lookup("mode").NoOptDefVal = bag.Raw // default
	}

	// aliases
	if *flag.BoolP("json", "j", false, "export in JSON format") {
		a.mode = bag.Json
	}

	if *flag.BoolP("jsonl", "J", false, "export in JSON Lines format") {
		a.mode = bag.Jsonl
	}

	if *flag.BoolP("xml", "X", false, "export in XML format") {
		a.mode = bag.Xml
	}

	if *flag.BoolP("sql", "S", false, "export in SQL format") {
		a.mode = bag.Sql
	}

	// standard options
	if *flag.BoolP("version", "v", false, "shows version") {
		fmt.Println(fox.Product, fox.Version)
		os.Exit(0)
	}

	flag.Usage = func() {
		fmt.Printf(fox.Usage, fox.Version)
		os.Exit(0)
	}

	flag.Parse()

	a.args = flag.Args()

	// flag checks
	if *head && *tail {
		sys.Exit("head or tail")
	}

	if *x && len(*H) > 0 {
		sys.Exit("hex or sum")
	}

	if *x && len(*filters) > 0 {
		sys.Exit("hex or pattern")
	}

	if *x && counts.Lines > 0 {
		sys.Exit("hex needs bytes")
	}

	// file limits
	if *head {
		limits.Head = *counts
	}

	if *tail {
		limits.Tail = *counts
	}

	// stdin piped
	if sys.IsPiped(os.Stdout) {
		a.print = true
	}

	// output mode
	if *w {
		a.print = true
		a.output.mode = types.Stats
	}

	if *s > 0 {
		a.print = true
		a.output.mode = types.Strings
		a.output.value = *s
	}

	if len(*H) > 0 {
		a.print = true
		a.output.mode = types.Hash
		a.output.value = *H
	}

	// render mode
	if *x {
		a.render.mode = mode.Hex
		a.output.mode = types.Hex
	}

	if len(*filters) > 0 {
		a.output.mode = types.Grep
	}

	return
}

func main() {
	a := argsParse()

	sys.Setup()

	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			sys.DumpErr(err, debug.Stack())
		}

		sys.Log.Close()
	}()

	_ = os.Remove(sys.Dump)

	hs := heapset.New(a.args)
	defer hs.ThrowAway()

	if a.print {
		hs.Print(a.output.mode, a.output.value)
	} else {
		hi := history.New()
		defer hi.Close()

		bg := bag.New(a.file, a.key, a.mode)
		defer bg.Close()

		fx := ui.New(a.render.mode)
		defer fx.Close()

		fx.Run(hs, hi, bg)
	}
}
