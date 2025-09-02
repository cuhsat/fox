package evidence

import (
	"os"
	"os/user"
	"time"
)

type Writer interface {
	Open(file *os.File, old bool, title string)

	Begin()
	Flush()

	WriteMeta(meta Meta)
	WriteLine(nr, grp int, str string)
}

type Evidence struct {
	Meta  Meta   `json:"meta"`
	Lines []Line `json:"lines"`
}

type Meta struct {
	User     *user.User `json:"user"`
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Size     int64      `json:"size"`
	Hash     []byte     `json:"hash"`
	Filters  []string   `json:"filters"`
	Bagged   time.Time  `json:"bagged"`
	Modified time.Time  `json:"modified"`
}

type Line struct {
	Nr  int    `json:"nr"`
	Grp int    `json:"grp"`
	Str string `json:"str"`
}

func Utc(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
