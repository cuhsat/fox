package user

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cuhsat/fx/pkg/fx/sys"
)

func Config(name string) (bool, string) {
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
