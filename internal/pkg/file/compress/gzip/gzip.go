package gzip

import (
	"compress/gzip"
	"io"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0x1F, 0x8B, 0x08,
	})
}

func Deflate(path string) string {
	a := sys.OpenFile(path)
	defer a.Close()

	r, err := gzip.NewReader(a)

	if err != nil {
		sys.Error(err)
		return path
	}

	defer r.Close()

	t := sys.TempFile("gzip")
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
