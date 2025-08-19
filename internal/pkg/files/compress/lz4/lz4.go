package lz4

import (
	"io"

	"github.com/pierrec/lz4"

	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/file"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0x04, 0x22, 0x4D, 0x18,
	})
}

func Deflate(path string) string {
	a := sys.Open(path)
	defer a.Close()

	r := lz4.NewReader(a)

	t := file.New(path)
	defer t.Close()

	_, err := io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
