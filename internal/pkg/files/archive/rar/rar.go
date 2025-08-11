package rar

import (
	"fmt"
	"io"
	"strings"

	"github.com/nwaples/rardecode"

	"github.com/hiforensics/fox/internal/pkg/files"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types/file"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0x52, 0x61, 0x72, 0x21, 0x1A, 0x07,
	})
}

func Deflate(path, pass string) (i []*files.Item) {
	a := sys.Open(path)
	defer a.Close()

	r, err := rardecode.NewReader(a, pass)

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

		t := file.New(fmt.Sprintf("%s/%s", path, h.Name))

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
