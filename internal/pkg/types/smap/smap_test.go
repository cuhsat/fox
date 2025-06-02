package smap

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/edsrzf/mmap-go"
)

func Testdata(name string) (*os.File, *mmap.MMap, error) {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		return nil, nil, errors.New("error")
	}

	p := filepath.Join(filepath.Dir(c), "..", "..", "..", "..", "test", "testdata", "json", name)

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
		f, m, err := Testdata("5MB.json")

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

func BenchmarkRender(b *testing.B) {
	b.Run("Benchmark Render", func(b *testing.B) {
		f, m, err := Testdata("5MB.json")

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
			s.Render()
		}
	})
}

func BenchmarkIndent(b *testing.B) {
	b.Run("Benchmark Indent", func(b *testing.B) {
		f, m, err := Testdata("5MB.json")

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
		f, m, err := Testdata("5MB.jsonl")

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

func BenchmarkGrep(b *testing.B) {
	b.Run("Benchmark Grep", func(b *testing.B) {
		f, m, err := Testdata("5MB.json")

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
		r := regexp.MustCompile(".*")

		b.ResetTimer()

		for b.Loop() {
			s.Grep(r)
		}
	})
}
