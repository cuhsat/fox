package zlib

import (
	"io"

	"github.com/klauspost/compress/zlib"

	"github.com/hiforensics/fox/internal/pkg/files"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types/file"
)

func Detect(path string) bool {
	for _, m := range [][]byte{
		{0x78, 0x01}, // no compression
		{0x78, 0x5E}, // fast compression
		{0x78, 0x9C}, // default compression
		{0x78, 0xDA}, // best compression
	} {
		if files.HasMagic(path, 0, m) {
			return true
		}
	}

	return false
}

func Deflate(path string) string {
	a := sys.Open(path)
	defer a.Close()

	r, err := zlib.NewReader(a)

	if err != nil {
		sys.Error(err)
		return path
	}

	defer r.Close()

	t := file.New(path + ".tmp")
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
