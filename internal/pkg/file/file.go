package file

import (
	"bytes"
	"io"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

type Item struct {
	Path string
	Name string
}

func HasMagic(p string, o int, m []byte) bool {
	buf := make([]byte, o+len(m))

	f := sys.OpenFile(p)
	defer f.Close()

	fi, err := f.Stat()

	if err != nil {
		sys.Error(err)
		return false
	}

	if fi.Size() < int64(o+len(m)) {
		return false
	}

	_, err = io.ReadFull(f, buf)

	if err != nil {
		sys.Error(err)
		return false
	}

	return bytes.Equal(buf[o:], m)
}
