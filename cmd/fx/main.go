package main

import (
    "flag"
    "os"

    "github.com/cuhsat/fx/internal/app"
    "github.com/cuhsat/fx/internal/sys"
    "github.com/cuhsat/fx/internal/sys/files/bag"
    "github.com/cuhsat/fx/internal/sys/files/history"
    "github.com/cuhsat/fx/internal/sys/heapset"
    "github.com/cuhsat/fx/internal/sys/types"
    "github.com/cuhsat/fx/internal/sys/types/mode"
)

const (
    Version = "dev"
)

func usage() {
    sys.Usage("usage: fx [-r] [-h | -t] [-n # | -c #] [-x | -e PATTERN] [-o FILE] [-b FILE] [PATH ... | -]")
}

func main() {
    c := new(types.Counts)
    l := types.GetLimits()
    e := types.GetFilters()

    // flags
    m := mode.Default
    r := flag.Bool("r", false, "Raw mode")
    x := flag.Bool("x", false, "Hex mode")
    
    // limits
    h := flag.Bool("h", false, "Head limit")
    t := flag.Bool("t", false, "Tail limit")

    // outputs
    o := flag.String("o", "", "Output file")
    b := flag.String("b", "", "Evidence file")

    // counts
    flag.IntVar(&c.Lines, "n", 0, "Lines count")
    flag.IntVar(&c.Bytes, "c", 0, "Bytes count")

    // filters
    flag.Var(e, "e", "Pattern")

    // standards
    v := flag.Bool("version", false, "Version")

    flag.Usage = usage
    flag.Parse()

    a := flag.Args()

    if len(a) == 0 {
        a = append(a, ".")
    }

    if *h && *t {
        sys.Usage("head or tail")
    }

    if c.Lines > 0 && c.Bytes > 0 {
        sys.Usage("lines or bytes")
    }

    if !*x && c.Bytes > 0 {
        sys.Usage("bytes needs hex")
    }

    if *x && len(*e) > 0 {
        sys.Usage("hex or pattern")
    }

    if *v {
        sys.Print("fx", Version)
        os.Exit(0)
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

    if len(*o) > 0 {
        *r = true
    }

    hs := heapset.NewHeapSet(a)
    defer hs.ThrowAway()

    if sys.IsPiped(os.Stdout) || *r {
        hs.Print(*o, *x)
        os.Exit(0)
    }

    hi := history.NewHistory()
    defer hi.Close()

    bag := bag.NewBag(*b)
    defer bag.Close()

    app := app.NewApp(m)
    defer app.Close()

    app.Run(hs, hi, bag)
}
