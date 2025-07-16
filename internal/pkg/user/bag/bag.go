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
	"github.com/cuhsat/fox/internal/pkg/arg"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
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

func New(args arg.ArgsBag) *Bag {
	var w writer

	if len(args.Path) == 0 {
		args.Path = arg.Bag
	}

	switch strings.ToLower(args.Mode) {
	case arg.Jsonl:
		w = NewJsonWriter(false)
		args.Path += ".jsonl"
	case arg.Json:
		w = NewJsonWriter(true)
		args.Path += ".json"
	case arg.Xml:
		w = NewXmlWriter()
		args.Path += ".xml"
	case arg.Sql:
		w = NewSqlWriter()
		args.Path += ".sqlite3"
	default:
		w = NewRawWriter()
		args.Path += ".txt"
	}

	return &Bag{
		Path: args.Path,
		key:  args.Key,
		file: nil,
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

	for _, str := range *h.FMap() {
		bag.w.SetLine(str.Nr, str.Str)
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
