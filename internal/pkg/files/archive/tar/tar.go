package tar

import (
	"archive/tar"
	"io"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
)

func Detect(path string) bool {
	return files.HasMagic(path, 257, []byte{
		0x75, 0x73, 0x74, 0x61, 0x72,
	})
}

func Deflate(path, _ string) (i []*files.Item) {
	a := fs.Open(path)
	defer a.Close()

	r := tar.NewReader(a)

	for {
		h, err := r.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			sys.Error(err)
			break
		}

		if strings.HasSuffix(h.Name, "/") {
			continue
		}

		t := fs.Create(filepath.Join(path, h.Name))

		_, err = io.Copy(t, r)
		_ = t.Close()

		if err != nil {
			sys.Error(err)
			continue
		}

		i = append(i, &files.Item{
			Path: t.Name(),
			Name: h.Name,
		})
	}

	return
}
