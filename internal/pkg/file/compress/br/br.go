package br

import (
	"io"

	"github.com/andybalholm/brotli"

	"github.com/hiforensics/fox/internal/pkg/file"
	"github.com/hiforensics/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0xCE, 0xB2, 0xCF, 0x81,
	})
}

func Deflate(path string) string {
	a := sys.OpenFile(path)
	defer a.Close()

	r := brotli.NewReader(a)

	t := sys.TempFile(path)
	defer t.Close()

	_, err := io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
