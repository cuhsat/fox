package heap

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "hash"
    "io"
    "strings"

    "github.com/cuhsat/cu/pkg/fs"
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
        fs.Panic("hash not supported")
    }

    _, err := h.file.Seek(0, io.SeekStart)

    if err != nil {
        fs.Panic(err)
    }

    _, err = io.Copy(a, h.file)
    
    if err != nil {
        fs.Panic(err)
    }

    b = a.Sum(nil)

    h.hash[algo] = b

    return b
}
