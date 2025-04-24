package heap

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "hash"
    "io"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
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
    sum, ok := h.hash[algo]

    if ok {
        return sum
    }

    var imp hash.Hash

    switch strings.ToLower(algo) {
    case Md5:
        imp = md5.New()
    case Sha1:
        imp = sha1.New()
    case Sha256:
        imp = sha256.New()
    default:
        fx.Error("hash not supported")

        return sum
    }

    f := fx.Open(h.Base)

    defer f.Close()

    _, err := io.Copy(imp, f)
    
    if err != nil {
        fx.Error(err)
    }

    sum = imp.Sum(nil)

    h.hash[algo] = sum

    return sum
}
