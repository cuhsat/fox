package file

import (
    "bytes"
    "io"

    "github.com/cuhsat/fx/internal/fx"
)

type Format func(string) []string

type Entry struct {
    Path string
    Name string
}

func HasMagic(p string, o int, m []byte) bool {
    buf := make([]byte, o + len(m))

    f := fx.Open(p)
    defer f.Close()

    fi, err := f.Stat()

    if err != nil {
        fx.Error(err)
        return false
    }

    if fi.Size() < int64(o + len(m)) {
        return false
    }

    _, err = io.ReadFull(f, buf)

    if err != nil {
        fx.Error(err)
        return false
    }

    return bytes.Equal(buf[o:], m)
}
