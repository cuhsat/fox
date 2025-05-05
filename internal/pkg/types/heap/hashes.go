package heap

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"strings"

	"github.com/cuhsat/fx/internal/pkg/sys"
)

const (
	Md5    = "md5"
	Sha1   = "sha1"
	Sha256 = "sha256"
)

type Hash map[string][]byte

func (h *Heap) Md5() ([]byte, error) {
	return h.Hashsum(Md5)
}

func (h *Heap) Sha1() ([]byte, error) {
	return h.Hashsum(Sha1)
}

func (h *Heap) Sha256() ([]byte, error) {
	return h.Hashsum(Sha256)
}

func (h *Heap) Hashsum(algo string) ([]byte, error) {
	algo = strings.ToLower(algo)

	h.RLock()
	sum, ok := h.hash[algo]
	h.RUnlock()

	if ok {
		return sum, nil
	}

	var imp hash.Hash

	switch algo {
	case Md5:
		imp = md5.New()
	case Sha1:
		imp = sha1.New()
	case Sha256:
		imp = sha256.New()
	default:
		return nil, errors.New("hash not supported")
	}

	f := sys.OpenFile(h.Base)

	defer f.Close()

	_, err := io.Copy(imp, f)

	if err != nil {
		return nil, err
	}

	sum = imp.Sum(nil)

	h.Lock()
	h.hash[algo] = sum
	h.Unlock()

	return sum, nil
}
