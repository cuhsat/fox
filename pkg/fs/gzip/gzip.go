package gzip

import (
    "bytes"
    "compress/gzip"
    "os"
    "io"

    "github.com/cuhsat/cu/pkg/fs"
)

var (
    Magic = [...]byte{0x1F, 0x8B, 0x08}
)

func Detect(p string) bool {
    r, err := os.Open(p)

    if err != nil {
        fs.Panic(err)
    }

    defer r.Close()

    var b [3]byte

    _, err = io.ReadFull(r, b[:])

    if err != nil {
        fs.Panic(err)
    }

    return bytes.Equal(b[:], Magic[:])
}

func Deflate(p string, f *os.File) string {
    z, err := os.Open(p)

    if err != nil {
        fs.Panic(err)
    }

    defer z.Close()

    r, err := gzip.NewReader(z)

    if err != nil {
        fs.Panic(err)
    }

    defer r.Close()

    _, err = io.Copy(f, r)

    if err != nil {
        fs.Panic(err)
    }

    defer f.Close()

    return f.Name()
}
