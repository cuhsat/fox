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

func BenchmarkMap(b *testing.B) {
	b.Run("Benchmark Map", func(b *testing.B) {
		f, m, err := testdata("test.txt")

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
		f, m, err := testdata("test.txt")

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
			s.Render(4)
		}
	})
}

func BenchmarkFormat(b *testing.B) {
	b.Run("Benchmark Format", func(b *testing.B) {
		f, m, err := testdata("test.json")

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
			s.Format(4)
		}
	})
}

func BenchmarkWrap(b *testing.B) {
	b.Run("Benchmark Wrap", func(b *testing.B) {
		f, m, err := testdata("test.txt")

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
			s.Wrap(4, 80)
		}
	})
}

func BenchmarkGrep(b *testing.B) {
	b.Run("Benchmark Grep", func(b *testing.B) {
		f, m, err := testdata("test.txt")

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

		re := regexp.MustCompile(".*")

		b.ResetTimer()

		for b.Loop() {
			s.Grep(re)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("Test Grep", func(t *testing.T) {
		f, m, err := testdata("test.txt")

		if err != nil {
			t.Fatal(err)
		}

		w, h := Map(m).Size()

		_ = m.Unmap()
		_ = f.Close()

		if w != 545 || h != 31107 {
			t.Fatal("wrong size")
		}
	})
}

func TestRender(t *testing.T) {
	t.Run("Test Render", func(t *testing.T) {
		b := []byte("\ttest\n")
		v := "    test\n"

		s := Map((*mmap.MMap)(&b)).Render(4)

		w, h := s.Size()

		if w != 8 || h != 1 {
			t.Fatal("wrong length")
		}

		if s.String() != v {
			t.Fatal("wrong string")
		}
	})
}

func TestFormat(t *testing.T) {
	t.Run("Test Format", func(t *testing.T) {
		b := []byte(`[{"test":123}]`)
		v := "[\n  {\n    \"test\": 123\n  }\n]\n"

		s := Map((*mmap.MMap)(&b)).Format(4)

		w, h := s.Size()

		if w != 15 || h != 5 {
			t.Fatal("wrong length")
		}

		if s.String() != v {
			t.Fatal("wrong string")
		}
	})
}

func TestWrap(t *testing.T) {
	t.Run("Test Wrap", func(t *testing.T) {
		b := []byte(`testtest`)
		v := "test\ntest\n"

		s := Map((*mmap.MMap)(&b)).Wrap(4, 4)

		w, h := s.Size()

		if w != 4 || h != 2 {
			t.Fatal("wrong length")
		}

		if s.String() != v {
			t.Fatal("wrong string")
		}
	})
}

func TestGrep(t *testing.T) {
	t.Run("Test Map", func(t *testing.T) {
		f, m, err := testdata("test.ioc")
		v := "test@example.org\nhttps://example.org\n"

		if err != nil {
			t.Fatal(err)
		}

		re := regexp.MustCompile("example")

		s := Map(m).Grep(re)

		_ = m.Unmap()
		_ = f.Close()

		if len(*s) != 2 {
			t.Fatal("wrong length")
		}

		if s.String() != v {
			t.Fatal("wrong string")
		}
	})
}

func testdata(name string) (*os.File, *mmap.MMap, error) {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		return nil, nil, errors.New("error")
	}

	p := filepath.Join(filepath.Dir(c), "..", "..", "..", "..", "testdata", name)

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
