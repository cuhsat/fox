package zip

import (
    "archive/zip"
    "bytes"
    "io"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/types"
)

var (
    magic = [...]byte{0x50, 0x4B, 0x03, 0x04}
)

func Detect(path string) bool {
    var buf [4]byte

    a := fx.Open(path)
    defer a.Close()

    fi, err := a.Stat()

    if err != nil {
        fx.Error(err)
        return false
    }

    if fi.Size() < 4 {
        return false
    }

    _, err = io.ReadFull(a, buf[:])

    if err != nil {
        fx.Error(err)
        return false
    }

    return bytes.Equal(buf[:], magic[:])
}

func Deflate(path string) (fe []*types.FileEntry) {
    r, err := zip.OpenReader(path)

    if err != nil {
        fx.Error(err)
 
        fe = append(fe, &types.FileEntry{
            Path: path,
            Name: path,
        })

        return
    }

    defer r.Close()

    for _, f := range r.File {
        if strings.HasSuffix(f.Name, "/") {
            continue
        }

        a, err := f.Open()

        if err != nil {
            fx.Error(err)
            continue
        }

        t := fx.Temp("zip", filepath.Ext(filepath.Base(f.Name)))

        _, err = io.Copy(t, a)

        t.Close()
        a.Close()

        if err != nil {
            fx.Error(err)
            continue
        }

        fe = append(fe, &types.FileEntry{
            Path: t.Name(),
            Name: f.Name,
        })
    }

    return
}
