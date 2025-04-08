package heap

import (
    "bytes"
    "io"
    "os"
    "runtime"
 
    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/smap"
    "github.com/cuhsat/cu/pkg/fs/limit"
    "github.com/edsrzf/mmap-go"
)

type Flag int

type Heap struct {
    Title string      // heap title
    Path  string      // file path
    Flag  Flag        // heap flags
    
    Limit limit.Limit // heap limit

    Head  int         // head offset
    Tail  int         // tail offset

    Chain []*Link     // filter chain

    MMap  mmap.MMap   // memory map
    SMap  smap.SMap   // string map current
    rmap  smap.SMap   // string map reserve

    hash  []byte      // file hash sum

    file  *os.File    // file handle
}

type Link struct {
    Name string    // filter name
    smap smap.SMap // filter string map
}

const (
    Normal Flag = iota
    StdIn
    StdOut
    StdErr
)

func (h *Heap) String() string {
    switch h.Flag {
    case StdIn:
        return "-"
    case StdOut:
        return h.Title
    case StdErr:
        return h.Title
    default:
        return h.Path
    }
}

func (h *Heap) Reload() {
    var err error

    h.ThrowAway()
    
    h.file, err = os.OpenFile(h.Path, os.O_RDONLY, 0644)

    if err != nil {
        fs.Panic(err)
    }

    h.MMap, err = mmap.Map(h.file, mmap.RDONLY, 0)

    if err != nil {
        fs.Panic(err)
    }

    // reduced mmap
    h.MMap, h.Head, h.Tail = h.Limit.ReduceMMap(h.MMap)

    // reduced smap
    h.SMap = h.Limit.ReduceSMap(smap.Map(h.MMap))

    h.rmap = h.SMap
    h.hash = h.hash[:0]

    h.ApplyFilter()
}

func (h *Heap) Length() int {
    return len(h.rmap)
}

func (h *Heap) Loaded() bool {
    return h.file != nil
}

func (h* Heap) Save() string {
    p := h.Path

    if h.Flag == StdOut || h.Flag == StdErr {
        p = h.String()
    }

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
        h.file = nil
    }

    runtime.GC()
}

func (h *Heap) write(w io.Writer) (n int, err error) {
    for i, s := range h.SMap {
        end := s.End

        if i < len(h.SMap)-1 {
            end += 1 // include breaks between strings
        }

        m, err := w.Write([]byte(h.MMap[s.Start:end]))

        if err != nil {
            return n, err
        }

        n += m
    }

    return n, nil
}
