package heap

import (
    "crypto/sha256"
    "io"

    "github.com/cuhsat/cu/pkg/fs"
)

func (h *Heap) Hash() []byte {
    if len(h.hash) == 0 {
        sha := sha256.New()

        _, err := io.Copy(sha, h.file)
        
        if err != nil {
            fs.Error(err)
        }

        h.hash = sha.Sum(nil)
    }

    return h.hash
}
