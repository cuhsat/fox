package limit

import (
    "github.com/cuhsat/cu/pkg/fs/smap"
    "github.com/edsrzf/mmap-go"
)

type Limit struct {
    Head Count // head limit
    Tail Count // tail limit
}

type Count struct {
    Lines int // lines count
    Bytes int // bytes count    
}

func (l *Limit) ReduceMMap(m mmap.MMap) (mmap.MMap, int, int) {
    h, t := 0, 0

    if l.Head.Bytes > 0 {
        h = min(l.Head.Bytes, len(m))
        m = m[:h]
    }

    if l.Tail.Bytes > 0 {
        t = max(len(m) - l.Tail.Bytes, 0)
        m = m[t:]
    }

    return m, h, t
}

func (l *Limit) ReduceSMap(s smap.SMap) smap.SMap {
    if l.Head.Lines > 0 {
        s = s[:min(l.Head.Lines, len(s))]
    }

    if l.Tail.Lines > 0 {
        s = s[max(len(s) - l.Tail.Lines, 0):]
    }

    return s
}
