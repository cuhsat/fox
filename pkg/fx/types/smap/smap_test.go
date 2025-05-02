package smap

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/edsrzf/mmap-go"
)

func Testdata(name string) string {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		return "error"
	}

	return filepath.Join(filepath.Dir(c), "..", "..", "..", "..", "test", "testdata", name)
}

func BenchmarkWrap(b *testing.B) {
	b.Run("Benchmark Wrap", func(b *testing.B) {
		f, err := os.OpenFile(Testdata("evtx.jsonl"), os.O_RDONLY, 0400)

		if err != nil {
			b.Fatal(err)
		}

		m, err := mmap.Map(f, mmap.RDONLY, 0)

		if err != nil {
			b.Fatal(err)
		}

		defer m.Unmap()

		s := Map(m)

		b.ResetTimer()

		for b.Loop() {
			s.Wrap(80)
		}
	})
}
