package main

import (
	"fmt"
	"os"
	"runtime/debug"

	flag "github.com/spf13/pflag"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/fox/ai"
	"github.com/cuhsat/fox/internal/fox/ui"
	"github.com/cuhsat/fox/internal/pkg/args"
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
	rm := mode.Default
	om := types.File
	var ov any

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
	limits := args.GetLimits()

	h := flag.BoolP("head", "h", false, "limit head of file by ...")
	t := flag.BoolP("tail", "t", false, "limit tail of file by ...")

	c := new(args.Counts)

	flag.IntVarP(&c.Lines, "lines", "n", 0, "number of lines")
	flag.IntVarP(&c.Bytes, "bytes", "c", 0, "number of bytes")

	if c.Lines == 0 {
		flag.Lookup("lines").NoOptDefVal = "10"
	}

	if c.Bytes == 0 {
		flag.Lookup("bytes").NoOptDefVal = "16"
	}

	// line filter
	filters := args.GetFilters()

	flag.VarP(filters, "regexp", "e", "filter for lines that matches pattern")

	// evidence bag
	f := flag.StringP("file", "f", "EVIDENCE", "file name of evidence bag")
	m := flag.StringP("mode", "m", "", "output mode")
	k := flag.StringP("key", "k", "", "key phrase for signing with HMAC")

	if len(*m) == 0 {
		flag.Lookup("mode").NoOptDefVal = bag.Text
	}

	// aliases
	j := flag.BoolP("json", "j", false, "export in JSON format")
	J := flag.BoolP("jsonl", "J", false, "export in JSON Lines format")
	M := flag.BoolP("markdown", "M", false, "export in Markdown format")
	X := flag.BoolP("xml", "X", false, "export in XML format")
	S := flag.BoolP("sql", "S", false, "export in SQL format")

	// standard options
	v := flag.BoolP("version", "v", false, "shows version")

	flag.Usage = func() {
		fmt.Printf(fox.Usage, fox.Version)
		os.Exit(0)
	}

	flag.Parse()

	a := flag.Args()

	// flag checks
	if *h && *t {
		sys.Exit("head or tail")
	}

	if *x && len(*H) > 0 {
		sys.Exit("hex or sum")
	}

	if *x && len(*filters) > 0 {
		sys.Exit("hex or pattern")
	}

	if *x && c.Lines > 0 {
		sys.Exit("hex needs bytes")
	}

	// features
	if *v {
		u, a := "no", "no"

		if ui.Build {
			u = "yes"
		}

		if ai.Build {
			a = "yes"
		}

		fmt.Println(fox.Product, fox.Version)
		fmt.Printf("Features UI: %s AI: %s\n", u, a)
		os.Exit(0)
	}

	// file limits
	if *h {
		limits.Head = *c
	}

	if *t {
		limits.Tail = *c
	}

	// evidence bag types
	if *j {
		*m = bag.Json
	}

	if *J {
		*m = bag.Jsonl
	}

	if *M {
		*m = bag.Markdown
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
		om = types.Stats
	}

	if *s > 0 {
		*p = true
		om = types.Strings
		ov = *s
	}

	if len(*H) > 0 {
		*p = true
		om = types.Hash
		ov = *H
	}

	// render mode
	if *x {
		rm = mode.Hex
		om = types.Hex
	}

	if len(*filters) > 0 {
		rm = mode.Grep
		om = types.Grep
	}

	sys.SetupLogger()

	_ = os.Remove(sys.FileDump)

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

	hs := heapset.New(a)
	defer hs.ThrowAway()

	if *p || !ui.Build {
		hs.Print(om, ov)
		return
	}

	hi := history.New()
	defer hi.Close()

	bg := bag.New(*f, *k, *m)
	defer bg.Close()

	u := ui.New(rm)
	defer u.Close()

	u.Run(hs, hi, bg)
}
