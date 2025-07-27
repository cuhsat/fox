package rar

import (
	"fmt"
	"io"
	"strings"

	"github.com/nwaples/rardecode"

	"github.com/hiforensics/fox/internal/pkg/file"
	"github.com/hiforensics/fox/internal/pkg/sys"
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0x52, 0x61, 0x72, 0x21, 0x1A, 0x07,
	})
}

func Deflate(path string) (i []*file.Item) {
	a := sys.OpenFile(path)
	defer a.Close()

	r, err := rardecode.NewReader(a, "")

	if err != nil {
		sys.Error(err)
		return
	}

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
