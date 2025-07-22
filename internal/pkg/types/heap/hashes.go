package heap

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"errors"
	"hash"
	"io"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/sys"

	"github.com/eciavatta/sdhash"
	"github.com/glaslos/ssdeep"
	"github.com/glaslos/tlsh"
)

const (
	Md5      = "md5"
	Sha1     = "sha1"
	Sha256   = "sha256"
	Sha3     = "sha3"
	Sha3_224 = "sha3-224"
	Sha3_256 = "sha3-256"
	Sha3_384 = "sha3-384"
	Sha3_512 = "sha3-512"
	Sdhash   = "sdhash"
	Ssdeep   = "ssdeep"
	Tlsh     = "tlsh"
)

type Hash map[string][]byte

func (h *Heap) Md5() ([]byte, error) {
	return h.HashSum(Md5)
}

func (h *Heap) Sha1() ([]byte, error) {
	return h.HashSum(Sha1)
}

func (h *Heap) Sha256() ([]byte, error) {
	return h.HashSum(Sha256)
}

func (h *Heap) Sha3() ([]byte, error) {
	return h.HashSum(Sha3)
}

func (h *Heap) Sdhash() ([]byte, error) {
	return h.HashSum(Sdhash)
}

func (h *Heap) Ssdeep() ([]byte, error) {
	return h.HashSum(Ssdeep)
}

func (h *Heap) Tlsh() ([]byte, error) {
	return h.HashSum(Tlsh)
}

func (h *Heap) HashSum(algo string) ([]byte, error) {
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
	case Sha3, Sha3_224:
		imp = sha3.New224()
	case Sha3_256:
		imp = sha3.New256()
	case Sha3_384:
		imp = sha3.New384()
	case Sha3_512:
		imp = sha3.New512()
	case Sdhash:
		imp = new(SDHash)
	case Ssdeep:
		imp = ssdeep.New()
	case Tlsh:
		imp = tlsh.New()
	default:
		return nil, errors.New("hash not supported")
	}

	imp.Reset()

	f := sys.OpenFile(h.Base)

	defer func(f sys.File) {
		_ = f.Close()
	}(f)

	b, err := io.ReadAll(f)

	if err != nil {
		return nil, err
	}

	_, err = imp.Write(b)

	if err != nil {
		return nil, err
	}

	sum = imp.Sum(nil)

	h.Lock()
	h.hash[algo] = sum
	h.Unlock()

	return sum, nil
}

type SDHash struct {
	sf   sdhash.SdbfFactory
	sdbf sdhash.Sdbf
}

func (s *SDHash) Reset() {
	s.sf = nil
	s.sdbf = nil
}

func (s *SDHash) BlockSize() int {
	return sdhash.BlockSize
}

func (s *SDHash) Size() int {
	return int(s.sdbf.Size())
}

func (s *SDHash) Sum(b []byte) []byte {
	s.sdbf = s.sf.Compute()

	return []byte(s.sdbf.String())
}

func (s *SDHash) Write(b []byte) (int, error) {
	var err error

	s.sf, err = sdhash.CreateSdbfFromBytes(b)

	return len(b), err
}
