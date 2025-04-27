package tar

import (
    "archive/tar"
    "io"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/pkg/fx/file"
    "github.com/cuhsat/fx/pkg/fx/sys"
)

func Detect(path string) bool {
    return file.HasMagic(path, 257, []byte{
        0x75, 0x73, 0x74, 0x61, 0x72,
    })
}

func Deflate(path string) (e []*file.Entry) {
    a := sys.Open(path)
    defer a.Close()

    r := tar.NewReader(a)

    for {
        h, err := r.Next()

        if err == io.EOF {
            break
        }

        if err != nil {
            sys.Error(err)
            break
        }

        if strings.HasSuffix(h.Name, "/") {
            continue
        }

        t := sys.Temp("tar", filepath.Ext(filepath.Base(h.Name)))

        _, err = io.Copy(t, r)

        t.Close()

        if err != nil {
            sys.Error(err)
            continue
        }

        e = append(e, &file.Entry{
            Path: t.Name(),
            Name: h.Name,
        })
    }

    return
}
