package bag

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/cuhsat/fox/api"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

type SqliteWriter struct {
	db    *sql.DB      // sql database
	entry *sqlEvidence // current entry
}

type sqlEvidence struct {
	created time.Time

	// file metadata
	file struct {
		path     string
		size     int64
		hash     string
		modified time.Time
		filters  []sqlValue
		lines    []sqlValue
	}

	// user metadata
	user struct {
		login string
		name  string
	}
}

type sqlValue struct {
	nr    int
	grp   int
	value string
}

func NewSqliteWriter() *SqliteWriter {
	return new(SqliteWriter)
}

func (w *SqliteWriter) Init(file *os.File, old bool, _ string) {
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

func (w *SqliteWriter) Start() {
	w.entry = new(sqlEvidence)
}

func (w *SqliteWriter) Flush() {
	tx, err := w.db.Begin()

	if err != nil {
		sys.Error(err)
		return
	}

	userId := w.insert(`users (login, name)`,
		w.entry.user.login,
		w.entry.user.name,
	)

	fileId := w.insert(`files (path, size, hash, modified)`,
		w.entry.file.path,
		w.entry.file.size,
		w.entry.file.hash,
		w.entry.file.modified,
	)

	for _, f := range w.entry.file.filters {
		_ = w.insert(`filters (file_id, nr, value)`,
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

func (w *SqliteWriter) WriteMeta(meta meta) {
	w.entry.created = meta.bagged.UTC()

	w.entry.file.path = meta.path
	w.entry.file.size = meta.size
	w.entry.file.modified = meta.modified.UTC()
	w.entry.file.hash = fmt.Sprintf("%x", meta.hash)

	for i, f := range meta.filters {
		w.entry.file.filters = append(w.entry.file.filters, sqlValue{
			nr: i, value: f,
		})
	}

	w.entry.user.login = meta.user.Username
	w.entry.user.name = meta.user.Name

}

func (w *SqliteWriter) WriteLine(nr, grp int, s string) {
	w.entry.file.lines = append(w.entry.file.lines, sqlValue{nr, grp, s})
}

func (w *SqliteWriter) insert(table string, v ...any) int64 {
	query := "INSERT INTO %s VALUES (%s);"

	return w.execute(fmt.Sprintf(query, table, fields(len(v))), v...)
}

func (w *SqliteWriter) execute(query string, v ...any) int64 {
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
