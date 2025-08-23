package bag

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	file "github.com/cuhsat/fox/internal/pkg/user"
)

type Bag struct {
	Path string        // bag path
	Mode flags.BagMode // bag mode

	file *os.File // file handle
	sign string   // sign phrase
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
	user     *user.User
	path     string
	size     int64
	hash     []byte
	filters  []string
	bagged   time.Time
	modified time.Time
}

func New() *Bag {
	var ws []writer

	bag := flags.Get().Bag

	if len(bag.Path) == 0 {
		bag.Path = flags.BagName
	}

	switch bag.Mode {
	case flags.BagModeNone:
	case flags.BagModeSqlite:
		ws = append(ws, NewSqliteWriter())
		bag.Path += ".sqlite3"
	case flags.BagModeJsonl:
		ws = append(ws, NewJsonWriter(false))
		bag.Path += ".jsonl"
	case flags.BagModeJson:
		ws = append(ws, NewJsonWriter(true))
		bag.Path += ".json"
	case flags.BagModeXml:
		ws = append(ws, NewXmlWriter())
		bag.Path += ".xml"
	case flags.BagModeText:
		ws = append(ws, NewTextWriter())
		bag.Path += ".bag"
	default:
		ws = append(ws, NewRawWriter())
		bag.Path += ".txt"
	}

	if len(bag.Url) > 0 {
		ws = append(ws, NewEcsWriter(bag.Url))
	}

	return &Bag{
		Path: bag.Path,
		Mode: bag.Mode,
		sign: bag.Sign,
		url:  bag.Url,
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

	usr, err := user.Current()

	if err != nil {
		sys.Error(err)
	}

	sum, err := h.HashSum(types.SHA256)

	if err != nil {
		sys.Error(err)
	}

	mod := time.Time.UTC(time.Now())

	fi, err := os.Stat(h.Base)

	if err == nil {
		mod = fi.ModTime()
	} else {
		sys.Error(err)
	}

	for _, w := range bag.ws {
		w.Start()

		w.WriteMeta(meta{
			user:     usr,
			path:     h.String(),
			size:     h.Len(),
			hash:     sum,
			filters:  h.Patterns(),
			bagged:   time.Now(),
			modified: mod,
		})

		for _, str := range *h.FMap() {
			w.WriteLine(str.Nr, str.Grp, str.Str)
		}

		w.Flush()
	}

	if bag.file != nil && len(bag.sign) > 0 {
		file.Sign(bag.Path, bag.sign)
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

	if bag.Mode != flags.BagModeNone {
		var err error

		bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

		if err != nil {
			sys.Error(err)
			return
		}
	}

	title := fmt.Sprintf("Forensic Examiner Evidence Bag (%s)", app.Version)

	for _, w := range bag.ws {
		w.Init(bag.file, old, title)
	}
}

func utc(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
