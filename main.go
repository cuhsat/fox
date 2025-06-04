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

func main() {
	// console
	rMode := mode.Default
	oMode := types.File
	var oValue any

	p := flag.BoolP("print", "p", false, "print to console (no UI)")
	x := flag.BoolP("hex", "x", false, "output file in hex / start in HEX mode")
	w := flag.BoolP("counts", "w", false, "output file line and byte counts")
	s := flag.IntP("strings", "s", 0, "output file ASCII and Unicode strings")
	H := flag.StringP("hash", "H", "", "output hash sum of file")

	if *s == 0 {
		flag.Lookup("strings").NoOptDefVal = "3"
	}

	if len(*H) == 0 {
		flag.Lookup("hash").NoOptDefVal = heap.Sha256
	}

	// file limits
	limits := types.GetLimits()

	head := flag.BoolP("head", "h", false, "limit head of file by ...")
	tail := flag.BoolP("tail", "t", false, "limit tail of file by ...")

	counts := new(types.Counts)

	flag.IntVarP(&counts.Lines, "lines", "n", 0, "number of lines")
	flag.IntVarP(&counts.Bytes, "bytes", "c", 0, "number of bytes")

	if counts.Lines == 0 {
		flag.Lookup("lines").NoOptDefVal = "10"
	}

	if counts.Bytes == 0 {
		flag.Lookup("bytes").NoOptDefVal = "16"
	}

	// line filter
	filters := types.GetFilters()

	flag.VarP(filters, "regexp", "e", "filter for lines that matches pattern")

	// evidence bag
	f := flag.StringP("file", "f", bag.Filename, "file name of evidence bag")
	m := flag.StringP("mode", "m", "", "output mode")
	k := flag.StringP("key", "k", "", "key phrase for signing with HMAC")

	if len(*m) == 0 {
		flag.Lookup("mode").NoOptDefVal = bag.Raw
	}

	// aliases
	j := flag.BoolP("json", "j", false, "export in JSON format")
	J := flag.BoolP("jsonl", "J", false, "export in JSON Lines format")
	X := flag.BoolP("xml", "X", false, "export in XML format")
	S := flag.BoolP("sql", "S", false, "export in SQL format")

	// standard options
	version := flag.BoolP("version", "v", false, "shows version")

	flag.Usage = func() {
		fmt.Printf(fox.Usage, fox.Version)
		os.Exit(0)
	}

	flag.Parse()

	args := flag.Args()

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

	// features
	if *version {
		fmt.Println(fox.Product, fox.Version)
		os.Exit(0)
	}

	// file limits
	if *head {
		limits.Head = *counts
	}

	if *tail {
		limits.Tail = *counts
	}

	// evidence bag types
	if *j {
		*m = bag.Json
	}

	if *J {
		*m = bag.Jsonl
	}

	if *X {
		*m = bag.Xml
	}

	if *S {
		*m = bag.Sql
	}

	// output mode
	if *w {
		*p = true
		oMode = types.Stats
	}

	if *s > 0 {
		*p = true
		oMode = types.Strings
		oValue = *s
	}

	if len(*H) > 0 {
		*p = true
		oMode = types.Hash
		oValue = *H
	}

	// render mode
	if *x {
		rMode = mode.Hex
		oMode = types.Hex
	}

	if len(*filters) > 0 {
		oMode = types.Grep
	}

	sys.Setup()

	_ = os.Remove(sys.Dump)

	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			sys.DumpErr(err, debug.Stack())
		}

		sys.Log.Close()
	}()

	if sys.IsPiped(os.Stdout) {
		*p = true
	}

	hs := heapset.New(args)
	defer hs.ThrowAway()

	if *p {
		hs.Print(oMode, oValue)
		return
	}

	hi := history.New()
	defer hi.Close()

	bg := bag.New(*f, *k, *m)
	defer bg.Close()

	u := ui.New(rMode)
	defer u.Close()

	u.Run(hs, hi, bg)
}
