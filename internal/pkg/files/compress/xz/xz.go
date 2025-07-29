package xz

import (
	"io"

	"github.com/ulikunitz/xz"

	"github.com/hiforensics/fox/internal/pkg/files"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types/file"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00,
	})
}

func Deflate(path string) string {
	a := sys.Open(path)
	defer a.Close()

	r, err := xz.NewReader(a)

	if err != nil {
		sys.Error(err)
		return path
	}

	t := file.New(path)
	defer t.Close()

	_, err = io.Copy(t, r)

	if err != nil {
		sys.Error(err)
		return path
	}

	return t.Name()
}
