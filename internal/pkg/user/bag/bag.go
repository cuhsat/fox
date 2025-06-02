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
	Raw   = "raw"
	Jsonl = "jsonl"
	Json  = "json"
	Xml   = "xml"
	Sql   = "sql"
)

const (
	filename = "EVIDENCE"
)

const (
	header = "Forensic Examiner Evidence Bag %s"
)

type Bag struct {
	Path string   // file path
	file *os.File // file handle
	key  string   // key phrase
	w    writer   // writer
}

type writer interface {
	Init(f *os.File, n bool, t string)
	Start()
	Finalize()
	WriteFile(p string, fs []string)
	WriteUser(u *user.User)
	WriteTime(t, f time.Time)
	WriteHash(b []byte)
	WriteLines(ns []int, ss []string)
}

func New(path, key, wt string) *Bag {
	var w writer
	var e string

	switch strings.ToLower(wt) {
	case Jsonl:
		w = NewJsonWriter(false)
		e = ".jsonl"
	case Json:
		w = NewJsonWriter(true)
		e = ".json"
	case Xml:
		w = NewXmlWriter()
		e = ".xml"
	case Sql:
		w = NewSqlWriter()
		e = ".db"
	case Raw:
		fallthrough
	default:
		w = NewRawWriter()
	}

	if len(path) == 0 {
		path = filename
	}

	if len(e) > 0 {
		path += e
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

	bag.w.WriteFile(h.String(), h.Patterns())
	bag.w.WriteUser(usr)
	bag.w.WriteTime(time.Now(), fi.ModTime())
	bag.w.WriteHash(sum)

	smap := *h.SMap()

	var ns []int
	var ss []string

	for _, s := range smap {
		ns = append(ns, s.Nr)
		ss = append(ss, s.Str)
	}

	bag.w.WriteLines(ns, ss)

	bag.w.Finalize()

	bag.sign()

	return true
}

func (bag *Bag) Close() {
	if bag.file == nil {
		_ = bag.file.Close()
	}
}

func (bag *Bag) init() bool {
	var err error

	is := sys.Exists(bag.Path)

	bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

	if err != nil {
		sys.Error(err)
		return false
	}

	bag.w.Init(bag.file, !is, fmt.Sprintf(header, fox.Version))

	return true
}

func (bag *Bag) sign() {
	var imp hash.Hash

	if len(bag.key) > 0 {
		imp = hmac.New(sha256.New, []byte(bag.key))
	} else {
		imp = sha256.New()
	}

	buf, err := os.ReadFile(bag.Path)

	if err != nil {
		sys.Error(err)
		return
	}

	imp.Write(buf)

	sum := fmt.Appendf(nil, "%x", imp.Sum(nil))

	err = os.WriteFile(bag.Path+".sha256", sum, 0600)

	if err != nil {
		sys.Error(err)
	}

	return
}

func writeln(f *os.File, s string) {
	_, err := fmt.Fprintln(f, s)

	if err != nil {
		sys.Error(err)
	}
}
