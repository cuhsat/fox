package files

import (
	"bytes"
	"io"

	"github.com/hiforensics/fox/internal/pkg/sys"
)

type Item struct {
	Path string
	Name string
}

type Deflate func(string, string) []*Item

func HasMagic(p string, o int, m []byte) bool {
	buf := make([]byte, o+len(m))

	f := sys.Open(p)
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
