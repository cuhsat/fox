package zstd

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0x28, 0xB5, 0x2F, 0xFD,
	})
}

func Deflate(path string) string {
	a := sys.OpenFile(path)
	defer a.Close()

	r, err := zstd.NewReader(a)

	if err != nil {
		sys.Error(err)
		return path
	}

	defer r.Close()

	b := strings.TrimSuffix(filepath.Base(path), ".zstd")

	t := sys.TempFile("zstd", filepath.Ext(b))
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
