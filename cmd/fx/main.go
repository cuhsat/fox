package main

import (
    "flag"
    "fmt"
    "os"
    "runtime/debug"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heapset"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/mode"
    "github.com/cuhsat/fx/internal/fx/user/bag"
    "github.com/cuhsat/fx/internal/fx/user/history"
    "github.com/cuhsat/fx/internal/ui"
)

const (
    Version = "dev"
)

func usage() {
    fmt.Println("usage: fx [-p] [-h | -t] [-n # | -c #] [-x | -e PATTERN] [-j | -J] [-o FILE] [PATH ... | -]")
    os.Exit(2)
}

func version() {
    fmt.Println("fx", version)
    os.Exit(0)
}

func main() {
    m := mode.Default

    c := new(types.Counts)
    l := types.GetLimits()
    e := types.GetFilters()

    // flags
    p := flag.Bool("p", false, "Print raw")
    x := flag.Bool("x", false, "Hexdump mode")
    j := flag.Bool("j", false, "JSON output")
    J := flag.Bool("J", false, "JSONL output")
    
    // limits
    h := flag.Bool("h", false, "Limit head")
    t := flag.Bool("t", false, "Limit tail")

    // output
    o := flag.String("o", "", "Evidence file")

    // counts
    flag.IntVar(&c.Lines, "n", 0, "Lines count")
    flag.IntVar(&c.Bytes, "c", 0, "Bytes count")

    // aliases
    flag.IntVar(&l.Head.Lines, "hn", 0, "Head lines count")
    flag.IntVar(&l.Head.Bytes, "hc", 0, "Head bytes count")
    flag.IntVar(&l.Tail.Lines, "tn", 0, "Tail lines count")
    flag.IntVar(&l.Tail.Bytes, "tc", 0, "Tail bytes count")

    // filters
    flag.Var(e, "e", "Pattern value")

    // standards
    v := flag.Bool("version", false, "Version")

    flag.Usage = usage
    flag.Parse()

    a := flag.Args()

    if len(a) == 0 {
        a = append(a, ".")
    }

    if *h && *t {
        fx.Exit("head or tail")
    }

    if c.Lines > 0 && c.Bytes > 0 {
        fx.Exit("lines or bytes")
    }

    if !*x && c.Bytes > 0 {
        fx.Exit("bytes needs hex")
    }

    if *x && len(*e) > 0 {
        fx.Exit("hex or pattern")
    }

    if *j && *J {
        fx.Exit("json or jsonl")
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

    if *x {
        m = mode.Hex
    }

    if len(*e) > 0 {
        m = mode.Grep
    }

    fx.SetupLogger()

    defer func() {
        if err := recover(); err != nil {
            fx.Dump(err, debug.Stack())
        }

        fx.Log.Close()
    }()

    if fx.IsPiped(os.Stdout) {
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

    bg := bag.New(*o, *j, *J)
    defer bg.Close()

    ui := ui.New(m)
    defer ui.Close()

    ui.Run(hs, hi, bg)
}
