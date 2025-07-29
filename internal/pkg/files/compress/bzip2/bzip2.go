package bzip2

import (
	"compress/bzip2"
	"io"

	"github.com/hiforensics/fox/internal/pkg/files"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types/file"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0x42, 0x5A, 0x68,
	})
}

func Deflate(path string) string {
	a := sys.Open(path)
	defer a.Close()

	r := bzip2.NewReader(a)

	t := file.New(path)
	defer t.Close()

	_, err := io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
