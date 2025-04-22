package gzip

import (
    "bytes"
    "compress/gzip"
    "io"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
)

var (
    magic = [...]byte{0x1F, 0x8B, 0x08}
)

func Detect(path string) bool {
    var buf [3]byte

    a := fx.Open(path)
    defer a.Close()

    fi, err := a.Stat()

    if err != nil {
        fx.Error(err)
        return false
    }

    if fi.Size() < 3 {
        return false
    }

    _, err = io.ReadFull(a, buf[:])

    if err != nil {
        fx.Error(err)
        return false
    }

    return bytes.Equal(buf[:], magic[:])
}

func Deflate(path string) string {
    a := fx.Open(path)
    defer a.Close()

    r, err := gzip.NewReader(a)

    if err != nil {
        fx.Error(err)
        return path
    }

    defer r.Close()

    b := strings.TrimSuffix(filepath.Base(path), ".gz")

    t := fx.Temp("gzip", filepath.Ext(b))
    defer t.Close()

    _, err = io.Copy(t, r)

    if err != nil {
        fx.Error(err)
        return path
    }

    return t.Name()
}
