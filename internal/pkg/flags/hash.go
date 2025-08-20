package flags

import (
	"errors"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/types"
)

type HashAlgo string

func (ha *HashAlgo) String() string {
	return strings.ToLower(string(*ha))
}

func (ha *HashAlgo) Type() string {
	return "HashAlgo"
}

func (ha *HashAlgo) Set(v string) error {
	switch strings.ToLower(v) {
	case types.CRC32IEEE:
		fallthrough
	case types.CRC64ISO:
		fallthrough
	case types.CRC64ECMA:
		fallthrough
	case types.SDHASH:
		fallthrough
	case types.SSDEEP:
		fallthrough
	case types.TLSH:
		fallthrough
	case types.MD5:
		fallthrough
	case types.SHA1:
		fallthrough
	case types.SHA256:
		fallthrough
	case types.SHA3:
		fallthrough
	case types.SHA3224:
		fallthrough
	case types.SHA3256:
		fallthrough
	case types.SHA3384:
		fallthrough
	case types.SHA3512:
		*ha = HashAlgo(v)
		return nil

	default:
		return errors.New("algorithm not recognized")
	}
}
