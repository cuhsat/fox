package zstd

import (
	"io"

	"github.com/klauspost/compress/zstd"

	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0x28, 0xB5, 0x2F, 0xFD,
	})
}

func Deflate(path string) string {
	a := fs.Open(path)
	defer a.Close()

	r, err := zstd.NewReader(a)

	if err != nil {
		sys.Error(err)
		return path
	}

	defer r.Close()

	t := fs.Create(path)
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
