package rar

import (
	"io"

	"github.com/gen2brain/go-unarr"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	for _, m := range [][]byte{
		{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00},       // v1.50
		{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x01, 0x00}, // v5.00
	} {
		if file.HasMagic(path, 0, m) {
			return true
		}
	}

	return false
}

func Deflate(path string) (i []*file.Item) {
	a, err := unarr.NewArchive(path)

	if err != nil {
		sys.Error(err)

		i = append(i, &file.Item{
			Path: path,
			Name: path,
		})

		return
	}

	defer a.Close()

	for {
		err := a.Entry()

		if err == io.EOF {
			break
		}

		if err != nil {
			sys.Error(err)
			break
		}

		t := sys.TempFile()

		_, err = io.Copy(t, a)
		_ = t.Close()

		if err != nil {
			sys.Error(err)
			continue
		}

		i = append(i, &file.Item{
			Path: t.Name(),
			Name: a.Name(),
		})
	}

	return
}
