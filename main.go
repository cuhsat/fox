package main

import (
	"fmt"
	"os"
	"runtime/debug"

	flag "github.com/spf13/pflag"

	"github.com/cuhsat/fx/internal/app/fx"
	"github.com/cuhsat/fx/internal/app/ui"
	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/heap"
	"github.com/cuhsat/fx/internal/pkg/types/heapset"
	"github.com/cuhsat/fx/internal/pkg/types/mode"
	"github.com/cuhsat/fx/internal/pkg/user/bag"
	"github.com/cuhsat/fx/internal/pkg/user/history"
)

const (
	Usage = ` _____                        _
|  ___|__  _ __ ___ _ __  ___(_) ___
| |_ / _ \| '__/ _ \ '_ \/ __| |/ __|
|  _| (_) | | |  __/ | | \__ \ | (__
|_|__\___/|_|  \___|_| |_|___/_|\___|
| ____|_  ____ _ _ __ ___ (_)_ __   ___ _ __
|  _| \ \/ / _' | '_ ' _ \| | '_ \ / _ \ '__|
| |___ >  < (_| | | | | | | | | | |  __/ |
|_____/_/\_\__,_|_| |_| |_|_|_| |_|\___|_| %s

The Swiss Army Knife for examining text files

Usage: fx [--print] [--count] [--hex] [--sum={md5|sha1|sha256}]
          [--head|tail] [--lines|bytes=NUMBER]
          [--regexp=PATTERN ...]
          [--file=FILE] [--key=KEY]
          [--mode={text,xml,json,jsonlines,markdown}] [-j|J|M|X]
          [PATH ...]

Positional arguments:
  PATH(s) to open or '-' for STDIN (default: current directory)

Console:
  -p, --print              print to console (no UI)
  -w, --count              output file line and byte counts
  -x, --hex                output file in hex / start in HEX mode
  -s, --sum={md5,sha1,sha256}      
                           output file hashsums (default: sha256)

File limits:
  -h, --head               limit head of file by ...
  -t, --tail               limit tail of file by ...
  -n, --lines=NUMBER       number of lines (default: 10)
  -c, --bytes=NUMBER       number of bytes (default: 16)

Line filter:
  -e, --regexp=PATTERN     filter for lines that matches pattern

Evidence bag:
  -f, --file=FILE          file name of evidence bag (default: "EVIDENCE")
  -m, --mode={text,xml,json,jsonl,markdown}
                           output mode (default: text)
  -k, --key=KEY            key phrase for signing with HMAC

Aliases:
  -j, --json               short for --mode=json
  -J, --jsonlines          short for --mode=jsonlines
  -M, --markdown           short for --mode=markdown
  -X, --xml                short for --mode=xml

Standard options:
      --help               shows this message
      --version            shows version

Full documentation: <https://github.com/cuhsat/fx/docs>
`
)

func main() {
	// console
	rm := mode.Default
	om := types.File

	p := flag.BoolP("print", "p", false, "print to console (no UI)")
	x := flag.BoolP("hex", "x", false, "output file in hex / start in HEX mode")
	w := flag.BoolP("count", "w", false, "output file line and byte count")
	s := flag.StringP("sum", "s", "", "output file hashsums")

	if len(*s) == 0 {
		flag.Lookup("sum").NoOptDefVal = heap.Sha256
	}

	// file limits
	l := types.Limits()

	h := flag.BoolP("head", "h", false, "limit head of file by ...")
	t := flag.BoolP("tail", "t", false, "limit tail of file by ...")

	c := new(types.Counts)

	flag.IntVarP(&c.Lines, "lines", "n", 0, "number of lines")
	flag.IntVarP(&c.Bytes, "bytes", "c", 0, "number of bytes")

	if c.Lines == 0 {
		flag.Lookup("lines").NoOptDefVal = "10"
	}

	if c.Bytes == 0 {
		flag.Lookup("bytes").NoOptDefVal = "16"
	}

	// line filter
	e := types.Filters()

	flag.VarP(e, "regexp", "e", "filter for lines that matches pattern")

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

	// standard options
	v := flag.BoolP("version", "v", false, "shows version")

	flag.Usage = func() {
		fmt.Printf(Usage, fx.Version)
		os.Exit(2)
	}

	flag.Parse()

	a := flag.Args()

	if len(a) == 0 {
		a = append(a, ".")
	}

	if *h && *t {
		sys.Exit("head or tail")
	}

	if *x && len(*s) > 0 {
		sys.Exit("hex or sum")
	}

	if *x && len(*e) > 0 {
		sys.Exit("hex or pattern")
	}

	if *x && c.Lines > 0 {
		sys.Exit("hex needs bytes")
	}

	if *v {
		fmt.Println(fx.Product, fx.Version)
		os.Exit(0)
	}

	if *h {
		l.Head = *c
	}

	if *t {
		l.Tail = *c
	}

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

	if *w {
		om = types.Count
	}

	if *x {
		rm = mode.Hex
		om = types.Hex
	}

	if len(*e) > 0 {
		rm = mode.Grep
		om = types.Grep
	}

	if len(*s) > 0 {
		om = types.Hash
	}

	sys.SetupLogger()

	os.Remove(sys.FileDump)

	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			sys.DumpErr(err, debug.Stack())
		}

		sys.Log.Close()
	}()

	if sys.IsPiped(os.Stdout) {
		*p = true
	}

	hs := heapset.New(a)
	defer hs.ThrowAway()

	if *p {
		hs.Print(om, *s)
		return
	}

	hi := history.New()
	defer hi.Close()

	bg := bag.New(*f, *k, *m)
	defer bg.Close()

	ui := ui.New(rm)
	defer ui.Close()

	ui.Run(hs, hi, bg)
}
