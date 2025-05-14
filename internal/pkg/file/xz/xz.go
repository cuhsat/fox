package xz

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00,
	})
}

func Deflate(path string) string {
	a := sys.OpenFile(path)
	defer a.Close()

	r, err := xz.NewReader(a)

	if err != nil {
		sys.Error(err)
		return path
	}

	b := strings.TrimSuffix(filepath.Base(path), ".xz")

	t := sys.TempFile("xz", filepath.Ext(b))
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
