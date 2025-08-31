package xz

import (
	"io"

	"github.com/ulikunitz/xz"

	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00,
	})
}

func Deflate(path string) string {
	a := fs.Open(path)
	defer a.Close()

	r, err := xz.NewReader(a)

	if err != nil {
		sys.Error(err)
		return path
	}

	t := fs.Create(path)
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
