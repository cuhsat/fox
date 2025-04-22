package tar

import (
    "archive/tar"
    "bytes"
    "io"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/types"
)

const (
    offset = 257
)

var (
    magic = [...]byte{0x75, 0x73, 0x74, 0x61, 0x72}
)

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

func Deflate(path string) (fe []*types.FileEntry) {
    a := fx.Open(path)
    defer a.Close()

    r := tar.NewReader(a)

    for {
        h, err := r.Next()

        if err == io.EOF {
            break
        }

        if err != nil {
            fx.Error(err)
            break
        }

        if strings.HasSuffix(h.Name, "/") {
            continue
        }

        t := fx.Temp("tar", filepath.Ext(filepath.Base(h.Name)))

        _, err = io.Copy(t, r)

        t.Close()

        if err != nil {
            fx.Error(err)
            continue
        }

        fe = append(fe, &types.FileEntry{
            Path: t.Name(),
            Name: h.Name,
        })
    }

    return
}
