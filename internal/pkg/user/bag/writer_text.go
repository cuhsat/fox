package bag

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

type TextWrite struct {
	file *os.File // file handle
}

func NewTextWriter() *TextWrite {
	return &TextWrite{
		file: nil,
	}
}

func (w *TextWrite) Init(f *os.File, old bool, title string) {
	w.file = f

	if old {
		return
	}

	_, err := fmt.Fprintf(w.file, "%s\n", title)

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) Start() {
	_, err := fmt.Fprintf(w.file, "\n")

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) Flush() {
	_, err := fmt.Fprintf(w.file, "\n")

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) SetFile(path string, size int64, fs []string) {
	_, err := fmt.Fprintf(w.file, "File: %s (%d bytes)\n", path, size)

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) SetUser(usr *user.User) {
	_, err := fmt.Fprintf(w.file, "User: %s (%s)\n", usr.Username, usr.Name)

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) SetTime(bag, mod time.Time) {
	_, err := fmt.Fprintf(w.file, "Time: %s / %s\n",
		bag.UTC().Format(time.RFC3339),
		mod.UTC().Format(time.RFC3339),
	)

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) SetHash(sum []byte) {
	_, err := fmt.Fprintf(w.file, "Hash: %x\n\n", sum)

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) SetLine(nr, grp int, s string) {
	_, err := fmt.Fprintf(w.file, "%d:%d:%s\n", nr, grp, s)

	if err != nil {
		sys.Error(err)
	}
}
