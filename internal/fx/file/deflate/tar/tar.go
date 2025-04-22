package tar

import (
    "archive/tar"
    "io"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/file"
)

func Detect(path string) bool {
    return file.HasMagic(path, 257, []byte{
        0x75, 0x73, 0x74, 0x61, 0x72,
    })
}

func Deflate(path string) (e []*file.Entry) {
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

        e = append(e, &file.Entry{
            Path: t.Name(),
            Name: h.Name,
        })
    }

    return
}
