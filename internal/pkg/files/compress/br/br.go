package br

import (
	"io"

	"github.com/andybalholm/brotli"

	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0xCE, 0xB2, 0xCF, 0x81,
	})
}

func Deflate(path string) string {
	a := fs.Open(path)
	defer a.Close()

	r := brotli.NewReader(a)

	t := fs.Create(path)
	defer t.Close()

	_, err := io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
