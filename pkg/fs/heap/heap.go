package heap

import (
    "bytes"
    "io"
    "os"
    "runtime"
 
    "github.com/cuhsat/cu/pkg/fs"
    "github.com/edsrzf/mmap-go"
)

type Heap struct {
    Path string    // file path
    Chain []*SLink // filter chain
    MMap mmap.MMap // memory map
    SMap SMap      // string map current
    rmap SMap      // string map reserve
    hash []byte    // file hash sum
    file *os.File  // file handle
}

type SLink struct {
    Name string // filter name
    smap SMap   // filter string map
}

func NewHeap(path string) *Heap {
    f, err := os.OpenFile(path, os.O_RDONLY, 0644)

    if err != nil {
        fs.Panic(err)
    }

    m, err := mmap.Map(f, mmap.RDONLY, 0)

    if err != nil {
        fs.Panic(err)
    }

    s := smap(m)

    return &Heap{
        Path: path,
        MMap: m,
        SMap: s,
        rmap: s,
        file: f,
    }
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

    h.SMap = smap(h.MMap)
    h.rmap = h.SMap
    h.hash = h.hash[:0]
}

func (h *Heap) Lines() int {
    return len(h.rmap)
}

func (h *Heap) Copy() []byte {
    var b bytes.Buffer

    err := h.lines(&b)

    if err != nil {
        fs.Panic(err)
    }

    return b.Bytes()
}

func (h* Heap) Save() string {
    fn := h.Path

    for _, l := range h.Chain {
        fn += "-" + l.Name
    }

    f, err := os.OpenFile(fn, fs.Override, 0644)

    if err != nil {
        fs.Panic(err)
    }

    defer f.Close()

    err = h.lines(f)

    if err != nil {
        fs.Panic(err)
    }

    return fn
}

func (h *Heap) AddFilter(value string) {
    h.SMap = h.filter([]byte(value))
    h.Chain = append(h.Chain, &SLink{
        Name: value,
        smap: h.SMap,
    })
}

func (h *Heap) DelFilter() {
    if len(h.Chain) > 0 {
        h.Chain = h.Chain[:len(h.Chain)-1]
    }

    if len(h.Chain) > 0 {
        h.SMap = h.Chain[len(h.Chain)-1].smap
    } else {
        h.SMap = h.rmap
    }
}

func (h *Heap) NoFilter() {
    h.Chain = h.Chain[:0]
    h.SMap = h.rmap
}

func (h *Heap) ThrowAway() {
    h.MMap.Unmap()
    h.file.Close()

    runtime.GC()
}

func (h *Heap) lines(w io.Writer) (err error) {
    for _, s := range h.SMap {
        _, err := w.Write([]byte(h.MMap[s.Start:s.End + 1]))

        if err != nil {
            return err
        }
    }

    return nil
}
