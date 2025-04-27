package main

import (
    "flag"
    "fmt"
    "os"
    "runtime/debug"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/args"
    "github.com/cuhsat/fx/internal/fx/heapset"
    "github.com/cuhsat/fx/internal/fx/sys"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/mode"
    "github.com/cuhsat/fx/internal/fx/user/bag"
    "github.com/cuhsat/fx/internal/fx/user/history"
    "github.com/cuhsat/fx/internal/ui"
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

usage: fx [-x] [-p] [-h|t] [-n|c #] [-e PATTERN]
          [-j|J] [-k KEY] [-o FILE]
          [-|PATH ...]

positional arguments:
  PATH to open (default: current dir)

mode:
  -x start in Hex mode

print:
  -p print to console (no UI)

limits:
  -h limit head of file by ...
  -t limit tail of file by ...
  -n # number of lines
  -c # number of bytes

filters:
  -e PATTERN to filter

evidence:
  -o FILE for evidence bag (default: EVIDENCE)
  -k KEY signing key phrase
  -j output in JSON format
  -J output in JSON lines format

options:
  --help    show help message
  --version show version info

`
)

func main() {
    e := types.Text
    m := mode.Default

    c := new(args.Counts)
    l := args.GetLimits()
    f := args.GetFilters()

    // flags
    p := flag.Bool("p", false, "")
    x := flag.Bool("x", false, "")
    j := flag.Bool("j", false, "")
    J := flag.Bool("J", false, "")
    
    // limits
    h := flag.Bool("h", false, "")
    t := flag.Bool("t", false, "")

    // output
    o := flag.String("o", "", "")
    k := flag.String("k", "", "")

    // counts
    flag.IntVar(&c.Lines, "n", 10, "")
    flag.IntVar(&c.Bytes, "c", 0, "")

    // filters
    flag.Var(f, "e", "Pattern value")

    // standards
    v := flag.Bool("version", false, "Version")

    flag.Usage = usage
    flag.Parse()

    a := flag.Args()

    if len(a) == 0 {
        a = append(a, ".")
    }

    if c.Bytes > 0 {
        c.Lines = 0
    }

    if *h && *t {
        sys.Exit("head or tail")
    }

    if *x && len(*f) > 0 {
        sys.Exit("hex or pattern")
    }

    if *x && c.Lines > 0 {
        sys.Exit("hex needs bytes")
    }

    if *v {
        version()
    }

    if *h {
        l.Head = *c
    }

    if *t {
        l.Tail = *c
    }

    if *j {
        e = types.Json
    }

    if *J {
        e = types.Jsonl
    }

    if *x {
        m = mode.Hex
    }

    if len(*f) > 0 {
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

    bg := bag.New(*o, *k, e)
    defer bg.Close()

    ui := ui.New(m)
    defer ui.Close()

    ui.Run(hs, hi, bg)
}

func usage() {
    fmt.Printf(Usage, fx.Version)
    os.Exit(2)
}

func version() {
    fmt.Println(fx.Product, fx.Version)
    os.Exit(0)
}
