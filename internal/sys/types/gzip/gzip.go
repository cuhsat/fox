package gzip

import (
    "bytes"
    "compress/gzip"
    "io"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/sys"
)

var (
    magic = [...]byte{0x1F, 0x8B, 0x08}
)

func Detect(p string) bool {
    g := sys.Open(p)

    defer g.Close()

    fi, err := g.Stat()

    if err != nil {
        sys.Fatal(err)
    }

    if fi.Size() < 3 {
        return false
    }

    var b [3]byte

    _, err = io.ReadFull(g, b[:])

    if err != nil {
        sys.Fatal(err)
    }

    return bytes.Equal(b[:], magic[:])
}

func Deflate(p string) string {
    g := sys.Open(p)

    defer g.Close()

    r, err := gzip.NewReader(g)

    if err != nil {
        sys.Fatal(err)
    }

    defer r.Close()

    n := filepath.Base(p)

    n = strings.TrimSuffix(n, ".Z")  // old style compress
    n = strings.TrimSuffix(n, ".gz") // new style compress

    f := sys.Temp("gzip", filepath.Ext(n))

    defer f.Close()

    _, err = io.Copy(f, r)

    if err != nil {
        sys.Fatal(err)
    }

    return f.Name()
}
