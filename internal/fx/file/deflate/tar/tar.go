package tar

import (
    "archive/tar"
    "bytes"
    "io"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
)

const (
    offset = 257
)

var (
    magic = [...]byte{0x75, 0x73, 0x74, 0x61, 0x72}
)

type TarEntry struct {
    Name, Path string
}

func Detect(path string) bool {
    var buf [offset+5]byte

    t := fx.Open(path)
    defer t.Close()

    fi, err := t.Stat()

    if err != nil {
        fx.Error(err)
        return false
    }

    if fi.Size() < offset+5 {
        return false
    }

    _, err = io.ReadFull(t, buf[:])

    if err != nil {
        fx.Error(err)
        return false
    }

    return bytes.Equal(buf[offset:], magic[:])
}

func Deflate(path string) (te []TarEntry) {
    r, err := tar.OpenReader(path)

    if err != nil {
        fx.Error(err)
 
        te = append(te, TarEntry{
            Name: path,
            Path: path,
        })

        return
    }

    defer r.Close()

    for _, f := range r.File {
        // if strings.HasSuffix(f.Name, "/") {
        //     continue
        // }

        a, err := f.Open()

        if err != nil {
            fx.Error(err)
            continue
        }

        b := filepath.Base(f.Name)

        t := fx.Temp("tar", filepath.Ext(b))

        _, err = io.Copy(t, a)

        t.Close()
        a.Close()

        if err != nil {
            fx.Error(err)
            continue
        }

        te = append(te, TarEntry{
            Name: f.Name,
            Path: t.Name(),
        })
    }

    return
}
