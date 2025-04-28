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

usage: fx [-x] [-p] [-h|t] [-n|c #] [-e PATTERN]
          [-j|J] [-k KEY] [-o FILE]
          [-|PATH ...]

positional arguments:
  PATH to open (default: current dir)

mode:
  --hex, -x        start in Hex mode

print:
  --print, -p      print to console (no UI)

limits:
  --head,  -h      limit head of file by ...
  --tail,  -t      limit tail of file by ...
  --lines, -n #    number of lines (default: 10)
  --bytes, -c #    number of bytes (default: 16)

filters:
  --regexp, -e     PATTERN to filter

evidence:
  --file,  -o      FILE for evidence bag (default: EVIDENCE)
  --key,   -k      KEY signing key phrase
  --json,  -j      output in JSON format
  --jsonl, -J      output in JSON lines format

options:
  --help           show help message
  --version        show version info

`
)

func main() {
    e := types.Text
    m := mode.Default

    c := new(types.Counts)
    l := types.Limits()
    f := types.Filters()

    // flags
    p := flag.BoolP("print", "p", false, "")
    x := flag.BoolP("hex", "x", false, "")
    j := flag.BoolP("json", "j", false, "")
    J := flag.BoolP("jsonl", "J", false, "")
    
    // limits
    h := flag.BoolP("head", "h", false, "")
    t := flag.BoolP("tail", "t", false, "")

    // output
    o := flag.StringP("file", "o", "", "")
    k := flag.StringP("key", "k", "", "")

    // counts
    flag.IntVarP(&c.Lines, "lines", "n", 0, "")
    flag.IntVarP(&c.Bytes, "bytes", "c", 0, "")

    if c.Lines == 0 {
        flag.Lookup("lines").NoOptDefVal = "10"
    }

    if c.Bytes == 0 {
        flag.Lookup("bytes").NoOptDefVal = "16"
    }

    // filters
    flag.VarP(f, "regexp", "e", "Pattern value")

    // standards
    v := flag.BoolP("version", "v", false, "Version")

    flag.Usage = usage
    flag.Parse()

    a := flag.Args()

    if len(a) == 0 {
        a = append(a, ".")
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
