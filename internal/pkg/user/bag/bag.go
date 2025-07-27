package bag

import (
	"fmt"
	"os"
	usr "os/user"
	"time"

	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/pkg/arg"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/user"
)

type Bag struct {
	Path string // bag path
	Mode string // bag mode

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

	switch args.Mode {
	case arg.BagNone:
	case arg.BagSqlite:
		ws = append(ws, NewSqliteWriter())
		args.Path += ".sqlite3"
	case arg.BagJsonl:
		ws = append(ws, NewJsonWriter(false))
		args.Path += ".jsonl"
	case arg.BagJson:
		ws = append(ws, NewJsonWriter(true))
		args.Path += ".json"
	case arg.BagXml:
		ws = append(ws, NewXmlWriter())
		args.Path += ".xml"
	case arg.BagText:
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
		Mode: args.Mode,
		key:  args.Key,
		url:  args.Url,
		file: nil,
		ws:   ws,
	}
}

func (bag *Bag) String() string {
	if bag.file != nil {
		return bag.Path
	} else {
		return bag.url
	}
}

func (bag *Bag) Put(h *heap.Heap) bool {
	bag.init()

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

	if bag.file != nil {
		user.Sign(bag.Path, bag.key)
	}

	return len(bag.ws) > 0
}

func (bag *Bag) Close() {
	if bag.file != nil {
		_ = bag.file.Close()
	}
}

func (bag *Bag) init() {
	old := sys.Exists(bag.Path)

	if bag.Mode != arg.BagNone {
		var err error

		bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

		if err != nil {
			sys.Error(err)
			return
		}
	}

	title := fmt.Sprintf("Forensic Examiner Evidence Bag %s", fox.Version)

	for _, w := range bag.ws {
		w.Init(bag.file, old, title)
	}
}
