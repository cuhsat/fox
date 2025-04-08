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

func (h *Heap) Hash(algo string) []byte {
    if len(h.hash) > 0 {
        return h.hash
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

    _, err := io.Copy(a, h.file)
    
    if err != nil {
        fs.Panic(err)
    }

    h.hash = a.Sum(nil)

    return h.hash
}
