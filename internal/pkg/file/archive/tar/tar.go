package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 257, []byte{
		0x75, 0x73, 0x74, 0x61, 0x72,
	})
}

func Deflate(path string) (i []*file.Item) {
	a := sys.OpenFile(path)
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

		t := sys.TempFile(fmt.Sprintf("%s/%s", path, h.Name))

		_, err = io.Copy(t, r)
		_ = t.Close()

		if err != nil {
			sys.Error(err)
			continue
		}

		i = append(i, &file.Item{
			Path: t.Name(),
			Name: h.Name,
		})
	}

	return
}
