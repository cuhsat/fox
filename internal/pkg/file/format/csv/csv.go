package csv

import (
	"path/filepath"
	"strings"

	"github.com/jfyne/csvd"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
)

func Detect(path string) bool {
	return filepath.Ext(strings.ToLower(path)) == ".csv"
}

func Format(path string) string {
	f := sys.OpenFile(path)
	defer f.Close()

	t := sys.TempFile("format")
	defer t.Close()

	r := csvd.NewReader(f)
	cols, err := r.ReadAll()

	if err != nil {
		sys.Error(err)
		return path
	}

	ls := make([]int, 0)

	// calculate row max length
	for _, rows := range cols {
		for i, row := range rows {
			if len(ls) < i+1 {
				ls = append(ls, 0)
			}

			ls[i] = max(text.Len(row), ls[i])
		}
	}

	var sb strings.Builder

	// prepad all rows
	for _, rows := range cols {
		for i, row := range rows {
			sb.WriteString(text.Pad(row, ls[i]+1))
		}

		sb.WriteRune('\n')

		_, err := t.WriteString(sb.String())

		if err != nil {
			sys.Error(err)
		}

		sb.Reset()
	}

	return t.Name()
}
