package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/cuhsat/fox/api"
	"github.com/cuhsat/fox/internal/pkg/files/evidence"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

const Ext = ".sqlite"

type Writer struct {
	db    *sql.DB // sql database
	entry *entry  // current entry
}

type entry struct {
	created time.Time

	// file metadata
	file struct {
		path     string
		size     int64
		hash     string
		modified time.Time
		filters  []value
		lines    []value
	}

	// user metadata
	user struct {
		login string
		name  string
	}
}

type value struct {
	nr    int
	grp   int
	value string
}

func New() *Writer {
	return new(Writer)
}

func (w *Writer) Open(file *os.File, old bool, _ string) {
	var err error

	_ = file.Close()

	w.db, err = sql.Open("sqlite", fmt.Sprintf("file:%s", file.Name()))

	if err != nil {
		sys.Panic(err)
	}

	// create the database from schema
	if !old {
		_, err = w.db.Exec(api.SchemaSql)

		if err != nil {
			sys.Panic(err)
		}
	}
}

func (w *Writer) Begin() {
	w.entry = new(entry)
}

func (w *Writer) Flush() {
	tx, err := w.db.Begin()

	if err != nil {
		sys.Error(err)
		return
	}

	userId := w.insert(`users (login, Name)`,
		w.entry.user.login,
		w.entry.user.name,
	)

	fileId := w.insert(`files (Path, Size, Hash, Modified)`,
		w.entry.file.path,
		w.entry.file.size,
		w.entry.file.hash,
		w.entry.file.modified,
	)

	for _, f := range w.entry.file.filters {
		_ = w.insert(`Filters (file_id, nr, value)`,
			fileId,
			f.nr,
			f.value,
		)
	}

	for _, l := range w.entry.file.lines {
		_ = w.insert(`lines (file_id, nr, grp, value)`,
			fileId,
			l.nr,
			l.grp,
			l.value,
		)
	}

	_ = w.insert(`evidence (user_id, file_id, created)`,
		userId,
		fileId,
		w.entry.created,
	)

	err = tx.Commit()

	if err != nil {
		sys.Error(err)
	}
}

func (w *Writer) WriteMeta(meta evidence.Meta) {
	w.entry.created = meta.Bagged.UTC()

	w.entry.file.path = meta.Path
	w.entry.file.size = meta.Size
	w.entry.file.modified = meta.Modified.UTC()
	w.entry.file.hash = fmt.Sprintf("%x", meta.Hash)

	for i, f := range meta.Filters {
		w.entry.file.filters = append(w.entry.file.filters, value{
			nr: i, value: f,
		})
	}

	w.entry.user.login = meta.User.Username
	w.entry.user.name = meta.User.Name
}

func (w *Writer) WriteLine(nr, grp int, str string) {
	w.entry.file.lines = append(w.entry.file.lines, value{nr, grp, str})
}

func (w *Writer) insert(table string, v ...any) int64 {
	query := "INSERT INTO %s VALUES (%s);"

	return w.execute(fmt.Sprintf(query, table, fields(len(v))), v...)
}

func (w *Writer) execute(query string, v ...any) int64 {
	res, err := w.db.Exec(query, v...)

	if err != nil {
		sys.Error(err)
		return 0
	}

	id, err := res.LastInsertId()

	if err != nil {
		sys.Error(err)
		return 0
	}

	return id
}

func fields(n int) string {
	var sb strings.Builder

	sb.WriteRune('?')

	for range n - 1 {
		sb.WriteString(", ?")
	}

	return sb.String()
}
