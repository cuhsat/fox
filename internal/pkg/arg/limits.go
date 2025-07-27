package arg

import (
	"github.com/edsrzf/mmap-go"

	"github.com/hiforensics/fox/internal/pkg/types/smap"
)

type Counts struct {
	Lines int // lines count
	Bytes int // bytes count
}

type Limits struct {
	Head Counts // head limit
	Tail Counts // tail limit
}

// singleton
var limits *Limits = nil

func GetLimits() *Limits {
	if limits == nil {
		limits = new(Limits)
	}

	return limits
}

func (l *Limits) ReduceMMap(m *mmap.MMap) *mmap.MMap {
	if l.Head.Bytes > 0 {
		r := make(mmap.MMap, min(l.Head.Bytes, len(*m)))
		copy(r, (*m)[:len(r)])
		return &r
	}

	if l.Tail.Bytes > 0 {
		r := make(mmap.MMap, min(len(*m), l.Tail.Bytes))
		copy(r, (*m)[max(len(*m)-len(r), 0):])
		return &r
	}

	return m
}

func (l *Limits) ReduceSMap(s *smap.SMap) *smap.SMap {
	if l.Head.Lines > 0 {
		r := make(smap.SMap, min(l.Head.Lines, len(*s)))
		copy(r, (*s)[:len(r)])
		return &r
	}

	if l.Tail.Lines > 0 {
		r := make(smap.SMap, min(len(*s), l.Tail.Lines))
		copy(r, (*s)[max(len(*s)-len(r), 0):])
		return &r
	}

	return s
}
