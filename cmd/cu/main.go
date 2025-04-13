package main

import (
    "flag"
    "os"
    "path/filepath"

    "github.com/cuhsat/cu/internal/app"
    "github.com/cuhsat/cu/internal/sys"
    "github.com/cuhsat/cu/internal/sys/files/history"
    "github.com/cuhsat/cu/internal/sys/heapset"
    "github.com/cuhsat/cu/internal/sys/types"
    "github.com/cuhsat/cu/internal/sys/types/mode"
)

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

    // output
    o := flag.String("o", "", "Output file")

    // counts
    flag.IntVar(&c.Lines, "n", 0, "Lines count")
    flag.IntVar(&c.Bytes, "c", 0, "Bytes count")

    // filters
    flag.Var(e, "e", "Pattern")

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

    app := app.NewApp(m)
    defer app.Close()

    app.Run(hs, hi)
}

func usage() {
    bin, err := os.Executable()

    if err != nil {
        sys.Fatal(err)
    }

    bin = filepath.Base(bin)

    sys.Usage("usage:", bin, "[-r] [-h | -t] [-n # | -c #] [-x | -e PATTERN] [-o FILE] [- | PATH ...]")
}
