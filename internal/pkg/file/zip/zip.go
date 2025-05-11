package zip

import (
	"archive/zip"
	"io"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0x50, 0x4B, 0x03, 0x04,
	})
}

func Deflate(path string) (i []*file.Item) {
	r, err := zip.OpenReader(path)

	if err != nil {
		sys.Error(err)

		i = append(i, &file.Item{
			Path: path,
			Name: path,
		})

		return
	}

	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		a, err := f.Open()

		if err != nil {
			sys.Error(err)
			continue
		}

		t := sys.TempFile("zip", filepath.Ext(filepath.Base(f.Name)))

		_, err = io.Copy(t, a)

		t.Close()
		a.Close()

		if err != nil {
			sys.Error(err)
			continue
		}

		i = append(i, &file.Item{
			Path: t.Name(),
			Name: f.Name,
		})
	}

	return
}
