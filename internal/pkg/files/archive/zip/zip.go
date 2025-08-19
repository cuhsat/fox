package zip

import (
	"fmt"
	"io"
	"strings"

	"github.com/cuhsat/zip/pkg/zip"

	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/file"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0x50, 0x4B, 0x03, 0x04,
	})
}

func Deflate(path, pass string) (i []*files.Item) {
	r, err := zip.OpenReader(path)

	if err != nil {
		sys.Error(err)

		i = append(i, &files.Item{
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

		if len(pass) > 0 {
			f.SetPassword(pass)
		}

		a, err := f.Open()

		if err != nil {
			sys.Error(err)
			continue
		}

		t := file.New(fmt.Sprintf("%s/%s", path, f.Name))

		_, err = io.Copy(t, a)

		_ = t.Close()
		_ = a.Close()

		if err != nil {
			sys.Error(err)
			continue
		}

		i = append(i, &files.Item{
			Path: t.Name(),
			Name: f.Name,
		})
	}

	return
}
