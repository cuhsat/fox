package limit

import (
    "github.com/cuhsat/cu/pkg/fs/smap"
    "github.com/edsrzf/mmap-go"
)

type Limit struct {
    head Count // head limit
    tail Count // tail limit
}

type Count struct {
    Lines int // lines count
    Bytes int // bytes count    
}

// singleton
var instance *Limit = nil

func GetLimit() *Limit {
    if instance == nil {
        instance = new(Limit);
    }

    return instance;
}

func SetHead(c Count) {
    GetLimit().head = c
}

func SetTail(c Count) {
    GetLimit().tail = c
}

func (l *Limit) ReduceMMap(m mmap.MMap) (mmap.MMap, int, int) {
    h, t := 0, 0

    if l.head.Bytes > 0 {
        h = min(l.head.Bytes, len(m))
        m = m[:h]
    }

    if l.tail.Bytes > 0 {
        t = max(len(m) - l.tail.Bytes, 0)
        m = m[t:]
    }

    return m, h, t
}

func (l *Limit) ReduceSMap(s smap.SMap) smap.SMap {
    if l.head.Lines > 0 {
        s = s[:min(l.head.Lines, len(s))]
    }

    if l.tail.Lines > 0 {
        s = s[max(len(s) - l.tail.Lines, 0):]
    }

    return s
}
