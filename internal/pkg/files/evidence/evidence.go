package evidence

import (
	"encoding/json"
	"os"
	"os/user"
	"time"
)

type Writer interface {
	Open(file *os.File, old bool, title string)

	Begin()
	Flush()

	WriteMeta(meta Meta)
	WriteLine(nr, grp int, s string)
}

type Evidence struct {
	Meta  Meta
	Lines []Line
}

type Meta struct {
	User     *user.User
	Name     string
	Path     string
	Size     int64
	Hash     []byte
	Filters  []string
	Bagged   time.Time
	Modified time.Time
}

type Line struct {
	Nr  int
	Grp int
	Str string
}

func New() *Evidence {
	return new(Evidence)
}

func (e *Evidence) String() string {
	buf, err := json.Marshal(e)

	if err == nil {
		return string(buf)
	} else {
		return err.Error()
	}
}

func (e *Evidence) SetMeta(meta Meta) {
	e.Meta = meta
}

func (e *Evidence) AddLine(nr, grp int, str string) {
	e.Lines = append(e.Lines, Line{nr, grp, str})
}

func Utc(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
