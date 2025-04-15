package text

import (
    "bytes"
    "io"
    "os"

    "github.com/cuhsat/fx/internal/sys"
)

func Lines(p string) (c int) {
    r, err := os.Open(p)

    if err != nil {
        sys.Fatal(err)
    }

    defer r.Close()
    
    b := make([]byte, 1024)

    for {
        n, err := r.Read(b)

        c += bytes.Count(b[:n], []byte{'\n'})

        switch err {
        case io.EOF:
            return
        case nil:
        default:
            sys.Fatal(err)
        }
    }

    return
}
