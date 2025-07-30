package arg

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/mode"
)

const (
	Bag = "evidence"
)

const (
	BagNone   = "none"
	BagText   = "text"
	BagJson   = "json"
	BagJsonl  = "jsonl"
	BagXml    = "xml"
	BagSqlite = "sqlite"
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
	Raw       bool
	NoConvert bool
	NoDeflate bool
	NoPlugins bool
}

type ArgsUI struct {
	Status string
	Mode   mode.Mode
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

	// actions
	H := flag.StringP("hash", "H", "", "")

	if len(*H) == 0 {
		flag.Lookup("hash").NoOptDefVal = "sha256" // default
	}

	s := flag.IntP("strings", "s", 0, "")

	if *s == 0 {
		flag.Lookup("strings").NoOptDefVal = "3" // default
	}

	flag.StringVarP(&args.Deflate, "deflate", "d", "", "")

	if len(args.Deflate) == 0 {
		flag.Lookup("deflate").NoOptDefVal = "-" // default
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

	N := flag.BoolP("no-bag", "", false, "")

	if len(args.Bag.Mode) == 0 {
		flag.Lookup("mode").NoOptDefVal = BagText // default
	}

	// disable
	flag.BoolVarP(&args.Opt.Raw, "raw", "r", false, "")
	flag.BoolVarP(&args.Opt.NoConvert, "no-convert", "", false, "")
	flag.BoolVarP(&args.Opt.NoDeflate, "no-deflate", "", false, "")
	flag.BoolVarP(&args.Opt.NoPlugins, "no-plugins", "", false, "")

	// aliases
	T := flag.BoolP("text", "T", false, "")
	j := flag.BoolP("json", "j", false, "")
	J := flag.BoolP("jsonl", "J", false, "")
	X := flag.BoolP("xml", "X", false, "")
	S := flag.BoolP("sqlite", "S", false, "")

	// interface
	flag.StringVarP(&args.UI.Status, "status", "", "", "")

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

	// checks
	if *head && *tail {
		sys.Exit("Can't specify both -h and -t")
	}

	if *x && len(*H) > 0 {
		sys.Exit("Can't specify both -x and -H")
	}

	if *x && *s > 0 {
		sys.Exit("Can't specify both -x and -s")
	}

	if *x && counts.Lines > 0 {
		sys.Exit("Can't specify both -x and -n")
	}

	if *x && len(filters.Patterns) > 0 {
		sys.Exit("Can't specify both -x and -e")
	}

	if len(args.Deflate) > 0 && args.Print.Active {
		sys.Exit("Can't specify both -d and -p")
	}

	if len(args.Deflate) > 0 && len(filters.Patterns) > 0 {
		sys.Exit("Can't specify both -d and -e")
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

	// disable
	if args.Opt.Raw {
		args.Opt.NoConvert = true
		args.Opt.NoDeflate = true
		args.Opt.NoPlugins = true
	}

	// aliases
	if *N {
		args.Bag.Mode = BagNone
	}

	if *T {
		args.Bag.Mode = BagText
	}

	if *j {
		args.Bag.Mode = BagJson
	}

	if *J {
		args.Bag.Mode = BagJsonl
	}

	if *X {
		args.Bag.Mode = BagXml
	}

	if *S {
		args.Bag.Mode = BagSqlite
	}

	if len(args.Bag.Mode) > 0 {
		args.Bag.Mode = strings.ToLower(args.Bag.Mode)
	}

	// interface
	if len(args.UI.Status) > 0 {
		re := regexp.MustCompile("[^-NWT]+")

		args.UI.Status = strings.ToUpper(args.UI.Status)
		args.UI.Status = re.ReplaceAllString(args.UI.Status, "")
	}

	// output mode
	if sys.Piped(os.Stdout) {
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
