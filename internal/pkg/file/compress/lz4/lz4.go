package lz4

import (
	"io"

	"github.com/pierrec/lz4"

	"github.com/hiforensics/fox/internal/pkg/file"
	"github.com/hiforensics/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0x04, 0x22, 0x4D, 0x18,
	})
}

func Deflate(path string) string {
	a := sys.OpenFile(path)
	defer a.Close()

	r := lz4.NewReader(a)

	t := sys.TempFile(path)
	defer t.Close()

	_, err := io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
