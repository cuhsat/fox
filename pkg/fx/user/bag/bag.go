package bag

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "hash"
    "os"
    "os/user"
    "time"

    "github.com/cuhsat/fx/pkg/fx/types"
    "github.com/cuhsat/fx/pkg/fx/types/heap"
    "github.com/cuhsat/fx/pkg/fx/sys"
)

const (
    filename = "EVIDENCE"
)

const (
    header = "FORENSIC EXAMINER EVIDENCE BAG"
)

type Bag struct {
    Path string   // file path
    file *os.File // file handle
    key string    // key phrase
    w writer      // writer
}

type writer interface {
    Init(f *os.File, n bool, t string)
    Start()
    Finalize()
    WriteFile(p string, f []string)
    WriteUser(u *user.User)
    WriteTime(t time.Time)
    WriteHash(b []byte)
    WriteLine(n int, s string)
}

func New(path, key string, ex types.Export) *Bag {
    var w writer
    var e string

    switch ex {
    case types.Jsonl:
        w = NewJsonWriter(false)
        e = ".jsonl"
    case types.Json:
        w = NewJsonWriter(true)
        e = ".json"
    default:
        w = NewTextWriter()
    }

    if len(path) == 0 {
        path = filename + e
    }

    return &Bag{
        Path: path,
        file: nil,
        key: key,
        w:   w,
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

    bag.w.Start()

    bag.w.WriteFile(h.String(), *types.Filters())
    bag.w.WriteUser(usr)
    bag.w.WriteTime(time.Now())
    bag.w.WriteHash(h.Sha256())

    for _, s := range h.SMap {
        bag.w.WriteLine(s.Nr, string(h.MMap[s.Start:s.End]))
    }

    bag.w.Finalize()

    bag.sign()

    return true
}

func (bag *Bag) Close() {
    if bag.file == nil {
        bag.file.Close()
    }
}

func (bag *Bag) init() bool {
    var err error

    is := sys.Exists(bag.Path)

    bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

    if err != nil {
        sys.Error(err)
        return false
    }

    bag.w.Init(bag.file, !is, header)

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

    sum := []byte(fmt.Sprintf("%x", imp.Sum(nil)))

    err = os.WriteFile(bag.Path + ".sha256", sum, 0600)

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
