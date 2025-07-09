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

// defaults
const (
	Bag = "evidence"
)

// options
const (
	Raw   = "raw"
	Json  = "json"
	Jsonl = "jsonl"
	Xml   = "xml"
	Sql   = "sql"
)

type Args struct {
	Args []string
	Dir  string

	Print ArgsPrint
	Bag   ArgsBag
	Opt   ArgsOpt
	UI    ArgsUI
}

type ArgsPrint struct {
	Active bool
	Mode   types.Print
	Value  any
}

type ArgsBag struct {
	Path string
	Mode string
	Key  string
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

	// console output
	flag.BoolVarP(&args.Print.Active, "print", "p", false, "print to console (no UI)")

	x := flag.BoolP("hex", "x", false, "output file in hex / start in HEX mode")
	w := flag.BoolP("counts", "w", false, "output file line and byte counts")
	s := flag.IntP("strings", "s", 0, "output file ASCII and Unicode strings")
	H := flag.StringP("hash", "H", "", "output hash sum of file")

	flag.StringVarP(&args.Dir, "output", "o", "", "output all files to folder")

	if *s == 0 {
		flag.Lookup("strings").NoOptDefVal = "3" // default
	}

	if len(*H) == 0 {
		flag.Lookup("hash").NoOptDefVal = "sha256" // default
	}

	if len(args.Dir) == 0 {
		flag.Lookup("output").NoOptDefVal = "out" // default
	}

	// file limits
	limits := GetLimits()

	head := flag.BoolP("head", "h", false, "limit head of file by ...")
	tail := flag.BoolP("tail", "t", false, "limit tail of file by ...")

	counts := new(Counts)

	flag.IntVarP(&counts.Lines, "lines", "n", 0, "number of lines")
	flag.IntVarP(&counts.Bytes, "bytes", "c", 0, "number of bytes")

	if counts.Lines == 0 {
		flag.Lookup("lines").NoOptDefVal = "10" // default
	}

	if counts.Bytes == 0 {
		flag.Lookup("bytes").NoOptDefVal = "16" // default
	}

	// line filter
	filters := GetFilters()

	flag.VarP(filters, "regexp", "e", "filter for lines that matches pattern")
	flag.IntVarP(&filters.Before, "before", "B", 0, "number of lines leading context before match")
	flag.IntVarP(&filters.After, "after", "A", 0, "number of lines trailing context after match")

	// evidence bag
	flag.StringVarP(&args.Bag.Path, "file", "f", Bag, "file name of evidence bag")
	flag.StringVarP(&args.Bag.Mode, "mode", "m", "", "output mode")
	flag.StringVarP(&args.Bag.Key, "key", "k", "", "key phrase for signing with HMAC")

	if len(args.Bag.Mode) == 0 {
		flag.Lookup("mode").NoOptDefVal = Raw // default
	}

	// aliases
	j := flag.BoolP("json", "j", false, "export in JSON format")
	J := flag.BoolP("jsonl", "J", false, "export in JSON Lines format")
	X := flag.BoolP("xml", "X", false, "export in XML format")
	S := flag.BoolP("sql", "S", false, "export in SQL format")

	// plugins
	flag.BoolVarP(&args.Opt.Skip, "skip", "a", false, "skip all automatic plugins")

	// standard options
	v := flag.BoolP("version", "v", false, "shows version")

	// show help
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

	if *&args.Print.Active && len(args.Dir) > 0 {
		sys.Exit("print or output")
	}

	// file limits
	if *head {
		limits.Head = *counts
	}

	if *tail {
		limits.Tail = *counts
	}

	// aliases
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
		args.Bag.Mode = Sql
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
