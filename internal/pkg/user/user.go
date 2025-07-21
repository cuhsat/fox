package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"os"
	"path/filepath"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

func File(name string) (bool, string) {
	dir, err := os.UserHomeDir()

	if err != nil {
		sys.Error(err)
		dir = "."
	}

	path := filepath.Join(dir, name)

	_, err = os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		return false, path
	} else if err != nil {
		sys.Error(err)
		return false, ""
	}

	return true, path
}

func Sign(path, key string) {
	var algo hash.Hash

	if len(path) == 0 {
		return
	}

	if len(key) > 0 {
		algo = hmac.New(sha256.New, []byte(key))
	} else {
		algo = sha256.New()
	}

	buf, err := os.ReadFile(path)

	if err != nil {
		sys.Error(err)
		return
	}

	algo.Write(buf)

	sum := fmt.Appendf(nil, "%x", algo.Sum(nil))

	err = os.WriteFile(path+".sha256", sum, 0600)

	if err != nil {
		sys.Error(err)
	}

	return
}
