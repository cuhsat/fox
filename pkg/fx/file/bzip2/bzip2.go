package bzip2

import (
	"compress/bzip2"
	"io"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fx/pkg/fx/file"
	"github.com/cuhsat/fx/pkg/fx/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0x42, 0x5A, 0x68,
	})
}

func Deflate(path string) string {
	a := sys.Open(path)
	defer a.Close()

	r := bzip2.NewReader(a)

	b := strings.TrimSuffix(filepath.Base(path), ".bz2")

	t := sys.Temp("bzip2", filepath.Ext(b))
	defer t.Close()

	_, err := io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
