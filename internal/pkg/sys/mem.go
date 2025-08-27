package sys

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"

	mem "github.com/cuhsat/memfile"
)

var fs = make(map[string]*mem.File)

type File = mem.Fileable

func Open(name string) File {
	mf, ok := fs[name]

	if ok {
		_, _ = mf.Seek(0, io.SeekStart)
		return mf // memory file
	}

	rf, err := os.OpenFile(name, os.O_RDONLY, 0400)

	if err == nil {
		return rf // regular file
	}

	Panic(err)
	return nil
}

func Create(name string) File {
	uniq := fmt.Sprintf("fox://%d/%s", rand.Uint64(), name)
	file := mem.New(uniq)

	fs[uniq] = file

	return file
}

func Exists(name string) bool {
	_, ok := fs[name]

	if ok {
		return true
	}

	_, err := os.Stat(name)

	return !errors.Is(err, os.ErrNotExist)
}

func Memory(name string) (File, bool) {
	f, ok := fs[name]

	return f, ok
}
