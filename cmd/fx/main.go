package main

import (
    "flag"
    "fmt"
    "os"
    "runtime/debug"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/files/bag"
    "github.com/cuhsat/fx/internal/fx/files/history"
    "github.com/cuhsat/fx/internal/fx/heapset"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/mode"
    "github.com/cuhsat/fx/internal/ui"
)

const (
    Version = "dev"
)

func main() {
    c := new(types.Counts)
    l := types.GetLimits()
    e := types.GetFilters()

    // flags
    m := mode.Default
    p := flag.Bool("p", false, "Print raw")
    x := flag.Bool("x", false, "Hexdump mode")
    
    // limits
    h := flag.Bool("h", false, "Limit head")
    t := flag.Bool("t", false, "Limit tail")

    // output
    o := flag.String("o", "", "Evidence file")

    // counts
    flag.IntVar(&c.Lines, "n", 0, "Lines count")
    flag.IntVar(&c.Bytes, "c", 0, "Bytes count")

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
        fx.Fatal("head or tail")
    }

    if c.Lines > 0 && c.Bytes > 0 {
        fx.Fatal("lines or bytes")
    }

    if !*x && c.Bytes > 0 {
        fx.Fatal("bytes needs hex")
    }

    if *x && len(*e) > 0 {
        fx.Fatal("hex or pattern")
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

    lf := fx.Init()
    defer lf.Close()

    if fx.IsPiped(os.Stdout) {
        *p = true
    }

    hs := heapset.New(a)
    defer hs.ThrowAway()

    if *p {
        hs.Print(*x)
        os.Exit(0)
    }

    hi := history.New()
    defer hi.Close()

    bg := bag.New(*o)
    defer bg.Close()

    ui := ui.New(m)
    defer ui.Close()

    defer func() {
        if err := recover(); err != nil {
            fx.Dump(err, debug.Stack())
        }
    }()
    
    ui.Run(hs, hi, bg)
}

func usage() {
    fmt.Println("usage: fx [-p] [-h | -t] [-n # | -c #] [-x | -e PATTERN] [-o FILE] [PATH ... | -]")
    os.Exit(2)
}

func version() {
    fmt.Println("fx", version)
    os.Exit(0)
}
