package types

import (
    "github.com/cuhsat/fx/pkg/fx/types/smap"
    "github.com/edsrzf/mmap-go"
)

type Counts struct {
    Lines int // lines count
    Bytes int // bytes count
}

type limits struct {
    Head Counts // head limit
    Tail Counts // tail limit
}

// singleton
var _limits *limits = nil

func Limits() *limits {
    if _limits == nil {
        _limits = new(limits);
    }

    return _limits;
}

func (l *limits) ReduceMMap(m mmap.MMap) (mmap.MMap, int, int) {
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

func (l *limits) ReduceSMap(s smap.SMap) smap.SMap {
    if l.Head.Lines > 0 {
        s = s[:min(l.Head.Lines, len(s))]
    }

    if l.Tail.Lines > 0 {
        s = s[max(len(s) - l.Tail.Lines, 0):]
    }

    return s
}
