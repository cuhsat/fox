package heap

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"errors"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"io"
	"strings"

	"github.com/eciavatta/sdhash"
	"github.com/glaslos/ssdeep"
	"github.com/glaslos/tlsh"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
)

type Hash map[string][]byte

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
	case types.CRC32IEEE:
		imp = crc32.NewIEEE()
	case types.CRC64ISO:
		imp = crc64.New(crc64.MakeTable(crc64.ISO))
	case types.CRC64ECMA:
		imp = crc64.New(crc64.MakeTable(crc64.ECMA))
	case types.SDHASH:
		imp = new(SDHash)
	case types.SSDEEP:
		imp = ssdeep.New()
	case types.TLSH:
		imp = tlsh.New()
	case types.MD5:
		imp = md5.New()
	case types.SHA1:
		imp = sha1.New()
	case types.SHA256:
		imp = sha256.New()
	case types.SHA3, types.SHA3224:
		imp = sha3.New224()
	case types.SHA3256:
		imp = sha3.New256()
	case types.SHA3384:
		imp = sha3.New384()
	case types.SHA3512:
		imp = sha3.New512()
	default:
		return nil, errors.New("hash not supported")
	}

	imp.Reset()

	var f sys.File

	// use deflate file not archive
	if h.Type == types.Deflate {
		f = sys.Open(h.Path)
	} else {
		f = sys.Open(h.Base)
	}

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
