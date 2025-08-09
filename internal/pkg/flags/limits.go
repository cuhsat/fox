package flags

import (
	"github.com/edsrzf/mmap-go"

	"github.com/hiforensics/fox/internal/pkg/types/smap"
)

type Limits struct {
	IsHead bool // is head limit
	IsTail bool // is tail limit
	Lines  int  // lines count
	Bytes  int  // bytes count
}

func (l *Limits) ReduceMMap(m *mmap.MMap) *mmap.MMap {
	if l.IsHead && l.Bytes > 0 {
		r := make(mmap.MMap, min(l.Bytes, len(*m)))
		copy(r, (*m)[:len(r)])
		return &r
	}

	if l.IsTail && l.Bytes > 0 {
		r := make(mmap.MMap, min(len(*m), l.Bytes))
		copy(r, (*m)[max(len(*m)-len(r), 0):])
		return &r
	}

	return m
}

func (l *Limits) ReduceSMap(s *smap.SMap) *smap.SMap {
	if l.IsHead && l.Lines > 0 {
		r := make(smap.SMap, min(l.Lines, len(*s)))
		copy(r, (*s)[:len(r)])
		return &r
	}

	if l.IsTail && l.Lines > 0 {
		r := make(smap.SMap, min(len(*s), l.Lines))
		copy(r, (*s)[max(len(*s)-len(r), 0):])
		return &r
	}

	return s
}
