package bag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

type EcsWriter struct {
	url   string    // export url
	title string    // export title
	entry *ecsEntry // current entry
}

type ecsEntry struct {
	Timestamp time.Time         `json:"@timestamp"`
	Message   string            `json:"message"`
	Labels    map[string]string `json:"labels"`

	Agent struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"agent"`

	Ecs struct {
		Version string `json:"version"`
	} `json:"ecs"`

	File struct {
		Mtime time.Time `json:"mtime"`
		Name  string    `json:"name"`
		Path  string    `json:"path"`
		Size  int64     `json:"size"`

		Hash struct {
			Sha256 string `json:"sha256"`
		} `json:"hash"`
	} `json:"file"`

	User struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"user"`
}

func NewEcsWriter(url string) *EcsWriter {
	return &EcsWriter{
		url: url,
	}
}

func (w *EcsWriter) Init(_ *os.File, _ bool, title string) {
	w.title = title
}

func (w *EcsWriter) Start() {
	w.entry = new(ecsEntry)
	w.entry.Labels = make(map[string]string)

	w.entry.Ecs.Version = "9.0.0"

	w.entry.Agent.Name = w.title
	w.entry.Agent.Type = app.Product
	w.entry.Agent.Version = app.Version[1:]
}

func (w *EcsWriter) Flush() {
	buf, err := json.Marshal(w.entry)

	if err != nil {
		sys.Error(err)
		return
	}

	res, err := http.Post(w.url, "application/json", bytes.NewBuffer(buf))

	if err != nil {
		sys.Error(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		sys.Error(http.StatusText(res.StatusCode))
	}
}

func (w *EcsWriter) WriteMeta(meta meta) {
	w.entry.Labels["filters"] = strings.Join(meta.filters, " > ")

	w.entry.Timestamp = meta.bagged.UTC()

	w.entry.File.Name = filepath.Base(meta.path)
	w.entry.File.Path = meta.path
	w.entry.File.Size = meta.size
	w.entry.File.Mtime = meta.modified.UTC()
	w.entry.File.Hash.Sha256 = fmt.Sprintf("%x", meta.hash)

	w.entry.User.Name = meta.user.Username
	w.entry.User.FullName = meta.user.Name
}

func (w *EcsWriter) WriteLine(nr, grp int, s string) {
	w.entry.Message += fmt.Sprintf("%d:%d: %s\n", nr, grp, s)
}
