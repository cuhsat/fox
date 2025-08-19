package cab

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-cabfile/cabfile"

	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/file"
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0x4D, 0x53, 0x43, 0x46,
	})
}

func Deflate(path, _ string) (i []*files.Item) {
	a := sys.Open(path)
	defer a.Close()

	r, err := cabfile.New(a)

	if err != nil {
		sys.Error(err)
		return
	}

	for _, s := range r.FileList() {
		if strings.HasSuffix(s, "/") {
			continue
		}

		h, err := r.Content(s)

		if err != nil {
			sys.Error(err)
			continue
		}

		t := file.New(fmt.Sprintf("%s/%s", path, s))

		_, err = io.Copy(t, h)
		_ = t.Close()

		if err != nil {
			sys.Error(err)
			continue
		}

		i = append(i, &files.Item{
			Path: t.Name(),
			Name: s,
		})
	}

	return
}
