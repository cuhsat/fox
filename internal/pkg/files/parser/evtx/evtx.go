package evtx

import (
	"github.com/0xrawsec/golang-evtx/evtx"

	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/file"
)

const (
	LF = 0xa
)

func Detect(path string) bool {
	return files.HasMagic(path, 0, []byte{
		0x45, 0x6C, 0x66, 0x46, 0x69, 0x6C, 0x65, 0x00,
	})
}

func Parse(path string) string {
	f := sys.Open(path)
	defer f.Close()

	r, err := evtx.New(f)
	defer r.Close()

	if err != nil {
		sys.Error(err)
		return path
	}

	t := file.New(path)
	defer t.Close()

	for e := range r.Events() {
		_, err := t.Write(evtx.ToJSON(e))

		if err != nil {
			sys.Error(err)
		}

		_, err = t.Write([]byte{LF})

		if err != nil {
			sys.Error(err)
		}
	}

	return t.Name()
}
