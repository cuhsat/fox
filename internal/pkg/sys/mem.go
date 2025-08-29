package sys

import (
	"os"

	"github.com/spf13/afero"
)

var fs = afero.NewMemMapFs()

type File = afero.File

func Open(path string) File {
	mf, err := fs.Open(path)

	if err == nil {
		return mf // memory file
	}

	rf, err := os.OpenFile(path, os.O_RDONLY, 0400)

	if err == nil {
		return rf // regular file
	}

	Panic(err)
	return nil
}

func Create(path string) File {
	f, err := fs.Create(path)

	if err != nil {
		Panic(err)
	}

	return f
}

func Exists(path string) bool {
	_, err := fs.Stat(path)

	if err == nil {
		return true
	}

	_, err = os.Stat(path)

	if err == nil {
		return true
	}

	return false
}
