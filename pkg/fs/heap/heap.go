package heap

import (
    "bytes"
    "io"
    "os"
    "runtime"
 
    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/smap"
    "github.com/edsrzf/mmap-go"
)

type Heap struct {
    Path  string    // file path
    Limit Limit     // line limit

    Chain []*Link   // filter chain

    MMap  mmap.MMap // memory map
    SMap  smap.SMap // string map current
    rmap  smap.SMap // string map reserve

    hash  []byte    // file hash sum

    file  *os.File  // file handle
}

type Limit struct {
    Head int // head count
    Tail int // tail count
}

type Link struct {
    Name string    // filter name
    smap smap.SMap // filter string map
}

func NewHeap(p string, l Limit) *Heap {
    h := Heap{
        Path: p,
        Limit: l,
    }

    h.Reload()
    
    return &h
}

func (h *Heap) String() string {
    if h.Path == fs.In {
        return "-"
    }

    return h.Path
}

func (h *Heap) Reload() {
    h.ThrowAway()
    var err error

    h.file, err = os.OpenFile(h.Path, os.O_RDONLY, 0644)

    if err != nil {
        fs.Panic(err)
    }

    h.MMap, err = mmap.Map(h.file, mmap.RDONLY, 0)

    if err != nil {
        fs.Panic(err)
    }

    h.SMap = smap.Map(h.MMap)

    l := len(h.SMap)

    if h.Limit.Head > 0 {
        h.SMap = h.SMap[:min(h.Limit.Head, l)]
    }

    if h.Limit.Tail > 0 {
        h.SMap = h.SMap[max(l-h.Limit.Tail, 0):]
    }

    h.rmap = h.SMap
    h.hash = h.hash[:0]
}

func (h *Heap) Length() int {
    return len(h.rmap)
}

func (h *Heap) Loaded() bool {
    return h.file != nil
}

func (h* Heap) Save() string {
    p := h.Path

    for _, l := range h.Chain {
        p += "-" + l.Name
    }

    f, err := os.OpenFile(p, fs.Override, 0644)

    if err != nil {
        fs.Panic(err)
    }

    defer f.Close()

    _, err = h.write(f)

    if err != nil {
        fs.Panic(err)
    }

    return p
}

func (h *Heap) Copy() []byte {
    var b bytes.Buffer

    _, err := h.write(&b)

    if err != nil {
        fs.Panic(err)
    }

    return b.Bytes()
}

func (h *Heap) ThrowAway() {
    if h.file != nil {
        h.MMap.Unmap()
        h.file.Close()
    }

    runtime.GC()
}

func (h *Heap) write(w io.Writer) (n int, err error) {
    for _, s := range h.SMap {
        m, err := w.Write([]byte(h.MMap[s.Start:s.End + 1]))

        if err != nil {
            return n, err
        }

        n += m
    }

    return n, nil
}
