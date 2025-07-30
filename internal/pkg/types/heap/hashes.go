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

	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types/file"

	"github.com/eciavatta/sdhash"
	"github.com/glaslos/ssdeep"
	"github.com/glaslos/tlsh"
)

const (
	MD5     = "md5"
	SHA1    = "sha1"
	SHA256  = "sha256"
	SHA3    = "sha3"
	SHA3224 = "sha3-224"
	SHA3256 = "sha3-256"
	SHA3384 = "sha3-384"
	SHA3512 = "sha3-512"
	SDHASH  = "sdhash"
	SSDEEP  = "ssdeep"
	TLSH    = "tlsh"
)

type Hash map[string][]byte

func (h *Heap) Md5() ([]byte, error) {
	return h.HashSum(MD5)
}

func (h *Heap) Sha1() ([]byte, error) {
	return h.HashSum(SHA1)
}

func (h *Heap) Sha256() ([]byte, error) {
	return h.HashSum(SHA256)
}

func (h *Heap) Sha3() ([]byte, error) {
	return h.HashSum(SHA3)
}

func (h *Heap) Sdhash() ([]byte, error) {
	return h.HashSum(SDHASH)
}

func (h *Heap) Ssdeep() ([]byte, error) {
	return h.HashSum(SSDEEP)
}

func (h *Heap) Tlsh() ([]byte, error) {
	return h.HashSum(TLSH)
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
	case MD5:
		imp = md5.New()
	case SHA1:
		imp = sha1.New()
	case SHA256:
		imp = sha256.New()
	case SHA3, SHA3224:
		imp = sha3.New224()
	case SHA3256:
		imp = sha3.New256()
	case SHA3384:
		imp = sha3.New384()
	case SHA3512:
		imp = sha3.New512()
	case SDHASH:
		imp = new(SDHash)
	case SSDEEP:
		imp = ssdeep.New()
	case TLSH:
		imp = tlsh.New()
	default:
		return nil, errors.New("hash not supported")
	}

	imp.Reset()

	f := sys.Open(h.Base)

	defer func(f file.File) {
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
	f sdhash.SdbfFactory
	s sdhash.Sdbf
}

func (sd *SDHash) Reset() {
	sd.f = nil
	sd.s = nil
}

func (sd *SDHash) BlockSize() int {
	return sdhash.BlockSize
}

func (sd *SDHash) Size() int {
	return int(sd.s.Size())
}

func (sd *SDHash) Sum(_ []byte) []byte {
	sd.s = sd.f.Compute()

	return []byte(strings.TrimRight(sd.s.String(), "\n"))
}

func (sd *SDHash) Write(b []byte) (int, error) {
	var err error

	sd.f, err = sdhash.CreateSdbfFromBytes(b)

	return len(b), err
}
