package types

import (
    "github.com/cuhsat/fx/internal/sys/types/smap"
    "github.com/edsrzf/mmap-go"
)

type Limits struct {
    Head Counts // head limit
    Tail Counts // tail limit
}

type Counts struct {
    Lines int // lines count
    Bytes int // bytes count    
}

// singleton instance
var limits *Limits = nil

func GetLimits() *Limits {
    if limits == nil {
        limits = new(Limits);
    }

    return limits;
}

func (l *Limits) MMapReduce(m mmap.MMap) (mmap.MMap, int, int) {
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

func (l *Limits) SMapReduce(s smap.SMap) smap.SMap {
    if l.Head.Lines > 0 {
        s = s[:min(l.Head.Lines, len(s))]
    }

    if l.Tail.Lines > 0 {
        s = s[max(len(s) - l.Tail.Lines, 0):]
    }

    return s
}
