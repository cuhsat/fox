package sevenzip

import (
	"io"

	"github.com/gen2brain/go-unarr"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C,
	})
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
