package zlib

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zlib"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	for _, m := range [][]byte{
		{0x78, 0x01}, // no compression
		{0x78, 0x5E}, // fast compression
		{0x78, 0x9C}, // default compression
		{0x78, 0xDA}, // best compression
	} {
		if file.HasMagic(path, 0, m) {
			return true
		}
	}

	return false
}

func Deflate(path string) string {
	a := sys.OpenFile(path)
	defer a.Close()

	r, err := zlib.NewReader(a)

	if err != nil {
		sys.Error(err)
		return path
	}

	defer r.Close()

	b := strings.TrimSuffix(filepath.Base(path), ".zz")

	t := sys.TempFile("zlib", filepath.Ext(b))
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
