package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"os"
	"path/filepath"

	"github.com/hiforensics/fox/internal/pkg/sys"
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
	var imp hash.Hash
	var ext string

	if len(path) == 0 {
		return
	}

	if len(key) > 0 && key != "-" {
		imp = hmac.New(sha256.New, []byte(key))
		ext = ".hmac_sha256"
	} else {
		imp = sha256.New()
		ext = ".sha256"
	}

	buf, err := os.ReadFile(path)

	if err != nil {
		sys.Error(err)
		return
	}

	imp.Write(buf)

	sum := fmt.Appendf(nil, "%x  %s\n", imp.Sum(nil), path)

	err = os.WriteFile(path+ext, sum, 0600)

	if err != nil {
		sys.Error(err)
	}

	return
}
