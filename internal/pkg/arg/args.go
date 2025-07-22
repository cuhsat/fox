package arg

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/mode"
)

const (
	Bag = "evidence"
)

const (
	Raw    = "raw"
	Text   = "text"
	Json   = "json"
	Jsonl  = "jsonl"
	Xml    = "xml"
	Sqlite = "sqlite"
)

type Args struct {
	Args []string

	Deflate string

	Print ArgsPrint
	Bag   ArgsBag
	Opt   ArgsOpt
	UI    ArgsUI
}

type ArgsPrint struct {
	Active bool
	NoFile bool
	NoLine bool
	Mode   types.Print
	Value  any
}

type ArgsBag struct {
	Path string
	Mode string
	Key  string
	Url  string
}

type ArgsOpt struct {
	Skip bool
}

type ArgsUI struct {
	Mode mode.Mode
}

// singleton
var args *Args = nil

func GetArgs() *Args {
	if args == nil {
		args = parse()
	}

	return args
}

func NewArgs() *Args {
	args := new(Args)

	args.Print.Mode = types.File
	args.UI.Mode = mode.Default

	return args
}

func parse() *Args {
	args := NewArgs()

	// console print
	flag.BoolVarP(&args.Print.Active, "print", "p", false, "")
	flag.BoolVarP(&args.Print.NoFile, "no-file", "", false, "")
	flag.BoolVarP(&args.Print.NoLine, "no-line", "", false, "")

	x := flag.BoolP("hex", "x", false, "")
	w := flag.BoolP("counts", "w", false, "")
	H := flag.StringP("hash", "H", "", "")

	if len(*H) == 0 {
		flag.Lookup("hash").NoOptDefVal = "sha256" // default
	}

	// carve strings
	s := flag.IntP("strings", "s", 0, "")

	if *s == 0 {
		flag.Lookup("strings").NoOptDefVal = "3" // default
	}

	// deflate file
	d := flag.StringP("deflate", "d", "", "")

	if len(*d) == 0 {
		flag.Lookup("deflate").NoOptDefVal = "out" // default
	}

	// file limits
	limits := GetLimits()

	head := flag.BoolP("head", "h", false, "")
	tail := flag.BoolP("tail", "t", false, "")

	counts := new(Counts)

	flag.IntVarP(&counts.Lines, "lines", "n", 0, "")
	flag.IntVarP(&counts.Bytes, "bytes", "c", 0, "")

	if counts.Lines == 0 {
		flag.Lookup("lines").NoOptDefVal = "10" // default
	}

	if counts.Bytes == 0 {
		flag.Lookup("bytes").NoOptDefVal = "16" // default
	}

	// line filter
	filters := GetFilters()

	flag.VarP(filters, "regexp", "e", "")

	C := flag.IntP("context", "C", 0, "")

	flag.IntVarP(&filters.Before, "before", "B", 0, "")
	flag.IntVarP(&filters.After, "after", "A", 0, "")

	// evidence bag
	flag.StringVarP(&args.Bag.Path, "file", "f", Bag, "")
	flag.StringVarP(&args.Bag.Mode, "mode", "m", "", "")
	flag.StringVarP(&args.Bag.Key, "key", "k", "", "")
	flag.StringVarP(&args.Bag.Url, "url", "u", "", "")

	if len(args.Bag.Mode) == 0 {
		flag.Lookup("mode").NoOptDefVal = Text // default
	}

	// aliases
	R := flag.BoolP("raw", "R", false, "")
	T := flag.BoolP("text", "T", false, "")
	j := flag.BoolP("json", "j", false, "")
	J := flag.BoolP("jsonl", "J", false, "")
	X := flag.BoolP("xml", "X", false, "")
	S := flag.BoolP("sqlite", "S", false, "")

	// plugins
	flag.BoolVarP(&args.Opt.Skip, "skip", "a", false, "")

	// standard options
	v := flag.BoolP("version", "v", false, "")

	// show usage
	flag.Usage = func() {
		fmt.Printf(fox.Usage, fox.Version)
		os.Exit(0)
	}

	flag.Parse()

	// show version
	if *v {
		fmt.Println(fox.Product, fox.Version)
		os.Exit(0)
	}

	args.Args = flag.Args()

	// flag checks
	if *head && *tail {
		sys.Exit("head or tail")
	}

	if *x && len(*H) > 0 {
		sys.Exit("hex or sum")
	}

	if *x && len(filters.Patterns) > 0 {
		sys.Exit("hex or pattern")
	}

	if *x && counts.Lines > 0 {
		sys.Exit("hex needs bytes")
	}

	if len(*d) > 0 && args.Print.Active {
		sys.Exit("deflate or print")
	}

	if len(*d) > 0 && len(filters.Patterns) > 0 {
		sys.Exit("deflate or pattern")
	}

	// file limits
	if *head {
		limits.Head = *counts
	}

	if *tail {
		limits.Tail = *counts
	}

	// line filter
	if *C > 0 {
		filters.Before = *C
		filters.After = *C
	}

	// aliases
	if *R {
		args.Bag.Mode = Raw
	}

	if *T {
		args.Bag.Mode = Text
	}

	if *j {
		args.Bag.Mode = Json
	}

	if *J {
		args.Bag.Mode = Jsonl
	}

	if *X {
		args.Bag.Mode = Xml
	}

	if *S {
		args.Bag.Mode = Sqlite
	}

	// deflate file
	if len(*d) > 0 {
		args.Deflate = *d
	}

	// output mode
	if sys.IsPiped(os.Stdout) {
		args.Print.Active = true
	}

	if len(filters.Patterns) > 0 {
		args.Print.Mode = types.Grep
	}

	if *w {
		args.Print.Active = true
		args.Print.Mode = types.Stats
	}

	if *s > 0 {
		args.Print.Active = true
		args.Print.Mode = types.Strings
		args.Print.Value = *s
	}

	if len(*H) > 0 {
		args.Print.Active = true
		args.Print.Mode = types.Hash
		args.Print.Value = *H
	}

	// ui mode
	if *x {
		args.UI.Mode = mode.Hex
		args.Print.Mode = types.Hex
	}

	return args
}
