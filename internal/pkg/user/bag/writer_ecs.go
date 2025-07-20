package bag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuhsat/fox/internal/fox"
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
	w.entry.Agent.Type = fox.Product
	w.entry.Agent.Version = fox.Version[1:]
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

func (w *EcsWriter) SetFile(path string, size int64, fs []string) {
	w.entry.Labels["filters"] = strings.Join(fs, " > ")

	w.entry.File.Name = filepath.Base(path)
	w.entry.File.Path = path
	w.entry.File.Size = size
}

func (w *EcsWriter) SetUser(usr *user.User) {
	w.entry.User.Name = usr.Username
	w.entry.User.FullName = usr.Name
}

func (w *EcsWriter) SetTime(bag, mod time.Time) {
	w.entry.Timestamp = bag.UTC()
	w.entry.File.Mtime = mod.UTC()
}

func (w *EcsWriter) SetHash(sum []byte) {
	w.entry.File.Hash.Sha256 = fmt.Sprintf("%x", sum)
}

func (w *EcsWriter) SetLine(nr, grp int, s string) {
	w.entry.Message += fmt.Sprintf("%d:%d: %s\n", nr, grp, s)
}
