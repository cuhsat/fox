package bag

import (
	"fmt"
	"os"
	usr "os/user"
	"strings"
	"time"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/pkg/arg"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/user"
)

type Bag struct {
	Path string // file path

	file *os.File // file handle
	key  string   // key phrase
	url  string   // url address
	ws   []writer // writers
}

type writer interface {
	Init(file *os.File, old bool, title string)

	Start()
	Flush()

	WriteMeta(meta meta)
	WriteLine(nr, grp int, s string)
}

type meta struct {
	user     *usr.User
	path     string
	size     int64
	hash     []byte
	filters  []string
	bagged   time.Time
	modified time.Time
}

func New(args arg.ArgsBag) *Bag {
	var ws []writer

	if len(args.Path) == 0 {
		args.Path = arg.Bag
	}

	switch strings.ToLower(args.Mode) {
	case arg.Sqlite:
		ws = append(ws, NewSqliteWriter())
		args.Path += ".sqlite3"
	case arg.Jsonl:
		ws = append(ws, NewJsonWriter(false))
		args.Path += ".jsonl"
	case arg.Json:
		ws = append(ws, NewJsonWriter(true))
		args.Path += ".json"
	case arg.Xml:
		ws = append(ws, NewXmlWriter())
		args.Path += ".xml"
	case arg.Text:
		ws = append(ws, NewTextWriter())
		args.Path += ".bag"
	default:
		ws = append(ws, NewRawWriter())
		args.Path += ".txt"
	}

	if len(args.Url) > 0 {
		ws = append(ws, NewEcsWriter(args.Url))
	}

	return &Bag{
		Path: args.Path,
		key:  args.Key,
		url:  args.Url,
		file: nil,
		ws:   ws,
	}
}

func (bag *Bag) Put(h *heap.Heap) bool {
	if bag.file == nil && !bag.init() {
		return false
	}

	u, err := usr.Current()

	if err != nil {
		sys.Error(err)
	}

	s, err := h.Sha256()

	if err != nil {
		sys.Error(err)
	}

	fi, err := os.Stat(h.Path)

	if err != nil {
		sys.Error(err)
	}

	for _, w := range bag.ws {
		w.Start()

		w.WriteMeta(meta{
			user:     u,
			path:     h.String(),
			size:     h.Len(),
			hash:     s,
			filters:  h.Patterns(),
			bagged:   time.Now(),
			modified: fi.ModTime(),
		})

		for _, str := range *h.FMap() {
			w.WriteLine(str.Nr, str.Grp, str.Str)
		}

		w.Flush()
	}

	user.Sign(bag.Path, bag.key)

	return true
}

func (bag *Bag) Close() {
	if bag.file != nil {
		_ = bag.file.Close()
	}
}

func (bag *Bag) init() bool {
	var err error

	old := sys.Exists(bag.Path)

	bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

	if err != nil {
		sys.Error(err)
		return false
	}

	title := fmt.Sprintf("Forensic Examiner Evidence Bag %s", fox.Version)

	for _, w := range bag.ws {
		w.Init(bag.file, old, title)
	}

	return true
}
