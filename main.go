package main

import (
	"fmt"
	"os"
	"runtime/debug"

	flag "github.com/spf13/pflag"

	"github.com/cuhsat/fx/pkg/fx"
	"github.com/cuhsat/fx/pkg/fx/sys"
	"github.com/cuhsat/fx/pkg/fx/types"
	"github.com/cuhsat/fx/pkg/fx/types/heapset"
	"github.com/cuhsat/fx/pkg/fx/types/mode"
	"github.com/cuhsat/fx/pkg/fx/user/bag"
	"github.com/cuhsat/fx/pkg/fx/user/history"
	"github.com/cuhsat/fx/pkg/ui"
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

usage: fx [--print] [--hex] [--head|tail] [--lines|bytes=NUM]
          [--json|jsonl] [--file=FILE] [--key=KEY]
          [--regexp=PATTERN ...]
          [-|PATH ...]

positional arguments:
  PATH(s) to open or '-' for STDIN (default: current directory)

general:
  -p, --print              print to console (no UI)
  -x, --hex                print or start in HEX mode

file limits:
  -h, --head               limit head of file by ...
  -t, --tail               limit tail of file by ...
  -n, --lines=NUM          number of lines (default: 10)
  -c, --bytes=NUM          number of bytes (default: 16)

line filter:
  -e, --regexp=PATTERN     filter for lines that matches pattern

evidence bag:
  -f, --file=FILE          file name of evidence bag (default: "EVIDENCE")
  -k, --key=KEY            key phrase for signing with HMAC
  -j, --json               export in JSON format
  -J, --jsonl              export in JSON Lines format

standard options:
      --help               shows this message
      --version            shows version

Full documentation: <https://github.com/cuhsat/fx/docs>
`
)

func main() {
	// general
	m := mode.Default

	p := flag.BoolP("print", "p", false, "print to console (no UI)")
	x := flag.BoolP("hex", "x", false, "print or start in HEX mode")

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
	fm := types.Text

	f := flag.StringP("file", "f", "EVIDENCE", "file name of evidence bag")
	k := flag.StringP("key", "k", "", "key phrase for signing with HMAC")
	j := flag.BoolP("json", "j", false, "export in JSON format")
	J := flag.BoolP("jsonl", "J", false, "export in JSON Lines format")

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
		fm = types.Json
	}

	if *J {
		fm = types.Jsonl
	}

	if *x {
		m = mode.Hex
	}

	if len(*e) > 0 {
		m = mode.Grep
	}

	sys.SetupLogger()

	os.Remove(sys.FileDump)

	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			sys.Dump(err, debug.Stack())
		}

		sys.Log.Close()
	}()

	if sys.IsPiped(os.Stdout) {
		*p = true
	}

	hs := heapset.New(a)
	defer hs.ThrowAway()

	if *p {
		hs.Print(*x)
		return
	}

	hi := history.New()
	defer hi.Close()

	bg := bag.New(*f, *k, fm)
	defer bg.Close()

	ui := ui.New(m)
	defer ui.Close()

	ui.Run(hs, hi, bg)
}
