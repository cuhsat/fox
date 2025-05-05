package gzip

import (
	"compress/gzip"
	"io"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fx/internal/pkg/file"
	"github.com/cuhsat/fx/internal/pkg/sys"
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

	b := strings.TrimSuffix(filepath.Base(path), ".gz")

	t := sys.TempFile("gzip", filepath.Ext(b))
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
