package bag

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

const (
	Filename = "evidence"
)

const (
	Raw   = "raw"
	Jsonl = "jsonl"
	Json  = "json"
	Xml   = "xml"
	Sql   = "sql"
)

type Bag struct {
	Path string // file path

	file *os.File // file handle
	key  string   // key phrase
	w    writer   // writer
}

type writer interface {
	Init(file *os.File, old bool, title string)

	Start()
	Flush()

	SetFile(path string, fs []string)
	SetUser(usr *user.User)
	SetTime(bag, mod time.Time)
	SetHash(sum []byte)
	SetLine(nr int, s string)
}

func New(path, key, mode string) *Bag {
	var w writer

	if len(path) == 0 {
		path = Filename
	}

	switch strings.ToLower(mode) {
	case Jsonl:
		w = NewJsonWriter(false)
		path += ".jsonl"
	case Json:
		w = NewJsonWriter(true)
		path += ".json"
	case Xml:
		w = NewXmlWriter()
		path += ".xml"
	case Sql:
		w = NewSqlWriter()
		path += ".db"
	default:
		w = NewRawWriter()
		path += ".txt"
	}

	return &Bag{
		Path: path,
		file: nil,
		key:  key,
		w:    w,
	}
}

func (bag *Bag) Put(h *heap.Heap) bool {
	if bag.file == nil && !bag.init() {
		return false
	}

	usr, err := user.Current()

	if err != nil {
		sys.Error(err)
	}

	sum, err := h.Sha256()

	if err != nil {
		sys.Error(err)
	}

	fi, err := os.Stat(h.Path)

	if err != nil {
		sys.Error(err)
	}

	bag.w.Start()

	bag.w.SetFile(h.String(), h.Patterns())
	bag.w.SetUser(usr)
	bag.w.SetTime(time.Now(), fi.ModTime())
	bag.w.SetHash(sum)

	for _, s := range *h.SMap() {
		bag.w.SetLine(s.Nr, s.Str)
	}

	bag.w.Flush()

	bag.hash()

	return true
}

func (bag *Bag) Close() {
	if bag.file == nil {
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

	bag.w.Init(bag.file, old, fmt.Sprintf("Forensic Examiner Evidence Bag %s", fox.Version))

	return true
}

func (bag *Bag) hash() {
	var algo hash.Hash

	if len(bag.key) > 0 {
		algo = hmac.New(sha256.New, []byte(bag.key))
	} else {
		algo = sha256.New()
	}

	buf, err := os.ReadFile(bag.Path)

	if err != nil {
		sys.Error(err)
		return
	}

	algo.Write(buf)

	sum := fmt.Appendf(nil, "%x", algo.Sum(nil))

	err = os.WriteFile(bag.Path+".sha256", sum, 0600)

	if err != nil {
		sys.Error(err)
	}

	return
}
