package archive

import (
	"fmt"
	"io"

	"github.com/gen2brain/go-unarr"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	for _, m := range [][]byte{
		{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}, // 7zip
		{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07}, // rar
		{0x50, 0x4B, 0x03, 0x04},             // zip
	} {
		if file.HasMagic(path, 0, m) {
			return true
		}
	}

	return file.HasMagic(path, 257, []byte{
		0x75, 0x73, 0x74, 0x61, 0x72, // tar
	})
}

func Deflate(path string) (i []*file.Item) {
	a, err := unarr.NewArchive(path)

	if err != nil {
		sys.Exit(err)

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
			sys.Exit(err)
		}

		b, err := a.ReadAll()

		if err != nil {
			sys.Exit(err)
		}

		t := sys.TempFile(fmt.Sprintf("%s/%s", path, a.Name()))

		_, err = t.Write(b)

		if err != nil {
			sys.Exit(err)
		}

		_ = t.Close()

		i = append(i, &file.Item{
			Path: t.Name(),
			Name: a.Name(),
		})
	}

	return
}
