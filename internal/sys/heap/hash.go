package heap

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "hash"
    "io"
    "os"
    "strings"

    "github.com/cuhsat/fx/internal/sys"
)

const (
    Md5    = "md5"
    Sha1   = "sha1"
    Sha256 = "sha256"
)

type Hash map[string][]byte

func (h *Heap) Md5() []byte {
    return h.HashSum(Md5)
}

func (h *Heap) Sha1() []byte {
    return h.HashSum(Sha1)
}

func (h *Heap) Sha256() []byte {
    return h.HashSum(Sha256)
}

func (h *Heap) HashSum(algo string) []byte {
    b, ok := h.hash[algo]

    if ok {
        return b
    }

    var a hash.Hash

    switch strings.ToLower(algo) {
    case Md5:
        a = md5.New()
    case Sha1:
        a = sha1.New()
    case Sha256:
        a = sha256.New()
    default:
        sys.Fatal("hash not supported")
    }

    f, err := os.Open(h.Base)

    if err != nil {
        sys.Fatal(err)
    }

    defer f.Close()

    _, err = io.Copy(a, f)
    
    if err != nil {
        sys.Fatal(err)
    }

    b = a.Sum(nil)

    h.hash[algo] = b

    return b
}
