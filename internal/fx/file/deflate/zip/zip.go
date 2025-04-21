package zip

import (
    "archive/zip"
    "bytes"
    "io"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
)

var (
    magic = [...]byte{0x50, 0x4B, 0x03, 0x04}
)

type ZipEntry struct {
    Name, Path string
}

func Detect(path string) bool {
    var buf [4]byte

    z := fx.Open(path)
    defer z.Close()

    fi, err := z.Stat()

    if err != nil {
        fx.Error(err)
        return false
    }

    if fi.Size() < 4 {
        return false
    }

    _, err = io.ReadFull(z, buf[:])

    if err != nil {
        fx.Error(err)
        return false
    }

    return bytes.Equal(buf[:], magic[:])
}

func Deflate(path string) (ze []ZipEntry) {
    r, err := zip.OpenReader(path)

    if err != nil {
        fx.Error(err)
 
        ze = append(ze, ZipEntry{
            Name: path,
            Path: path,
        })

        return
    }

    defer r.Close()

    for _, f := range r.File {
        if strings.HasSuffix(f.Name, "/") {
            continue
        }

        z, err := f.Open()

        if err != nil {
            fx.Error(err)
            continue
        }

        b := filepath.Base(f.Name)

        t := fx.Temp("zip", filepath.Ext(b))

        _, err = io.Copy(t, z)

        t.Close()
        z.Close()

        if err != nil {
            fx.Error(err)
            continue
        }

        ze = append(ze, ZipEntry{
            Name: f.Name,
            Path: t.Name(),
        })
    }

    return
}
