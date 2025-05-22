package smap

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/edsrzf/mmap-go"
)

func Testdata(name string) (*os.File, *mmap.MMap, error) {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		return nil, nil, errors.New("error")
	}

	p := filepath.Join(filepath.Dir(c), "..", "..", "..", "..", "test", "testdata", "jsonl", name)

	f, err := os.OpenFile(p, os.O_RDONLY, 0400)

	if err != nil {
		return nil, nil, err
	}

	m, err := mmap.Map(f, mmap.RDONLY, 0)

	if err != nil {
		return nil, nil, err
	}

	return f, &m, nil
}

func BenchmarkMap(b *testing.B) {
	b.Run("Benchmark Map", func(b *testing.B) {
		f, m, err := Testdata("evtx.jsonl")

		if err != nil {
			b.Fatal(err)
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		defer func(m *mmap.MMap) {
			_ = m.Unmap()
		}(m)

		b.ResetTimer()

		for b.Loop() {
			Map(m)
		}
	})
}

func BenchmarkIndent(b *testing.B) {
	b.Run("Benchmark Indent", func(b *testing.B) {
		f, m, err := Testdata("evtx.jsonl")

		if err != nil {
			b.Fatal(err)
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		defer func(m *mmap.MMap) {
			_ = m.Unmap()
		}(m)

		s := Map(m)

		b.ResetTimer()

		for b.Loop() {
			s.Indent()
		}
	})
}

func BenchmarkWrap(b *testing.B) {
	b.Run("Benchmark Wrap", func(b *testing.B) {
		f, m, err := Testdata("evtx.jsonl")

		if err != nil {
			b.Fatal(err)
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		defer func(m *mmap.MMap) {
			_ = m.Unmap()
		}(m)

		s := Map(m)

		b.ResetTimer()

		for b.Loop() {
			s.Wrap(80)
		}
	})
}
