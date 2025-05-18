package evtx

import (
	"github.com/0xrawsec/golang-evtx/evtx"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

const (
	LF byte = 0xa
)

func Detect(path string) bool {
	return file.HasMagic(path, 0, []byte{
		0x45, 0x6C, 0x66, 0x46, 0x69, 0x6C, 0x65, 0x00,
	})
}

func Parse(path string) string {
	f, err := evtx.OpenDirty(path)
	defer f.Close()

	if err != nil {
		sys.Error(err)
		return path
	}

	t := sys.TempFile()
	defer t.Close()

	for e := range f.Events() {
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
