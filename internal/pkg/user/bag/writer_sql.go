package bag

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"time"

	_ "embed"
	_ "modernc.org/sqlite"

	"github.com/cuhsat/fox/api"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

var (
	// go:embed bag.sql
	SqlSchema string
)

type SqlWriter struct {
	db    *sql.DB      // sql database
	entry *sqlEvidence // current entry
}

type sqlEvidence struct {
	created time.Time

	// user metadata
	user struct {
		login string
		name  string
	}

	// file metadata
	file struct {
		path     string
		hash     string
		modified time.Time
		filters  []sqlData
		lines    []sqlData
	}
}

type sqlData struct {
	nr    int
	value string
}

func NewSqlWriter() *SqlWriter {
	return new(SqlWriter)
}

func (w *SqlWriter) Init(f *os.File, n bool, _ string) {
	var err error

	f.Close()

	w.db, err = sql.Open("sqlite", fmt.Sprintf("file:%s", f.Name()))

	if err != nil {
		sys.Panic(err)
	}

	if !n {
		return
	}

	_, err = w.db.Exec(api.SqlSchema)

	if err != nil {
		sys.Panic(err)
	}
}

func (w *SqlWriter) Start() {
	w.entry = new(sqlEvidence)
}

func (w *SqlWriter) Finalize() {
	tx, err := w.db.Begin()

	if err != nil {
		sys.Error(err)
		return
	}

	user_id := w.insert(`users (login, name)`,
		w.entry.user.login,
		w.entry.user.name,
	)

	file_id := w.insert(`files (path, hash, modified)`,
		w.entry.file.path,
		w.entry.file.hash,
		w.entry.file.modified,
	)

	for _, f := range w.entry.file.filters {
		_ = w.insert(`filters (file_id, nr, value)`,
			file_id,
			f.nr,
			f.value,
		)
	}

	for _, l := range w.entry.file.lines {
		_ = w.insert(`lines (file_id, nr, value)`,
			file_id,
			l.nr,
			l.value,
		)
	}

	_ = w.insert(`evidence (user_id, file_id, created)`,
		user_id,
		file_id,
		w.entry.created,
	)

	err = tx.Commit()

	if err != nil {
		sys.Error(err)
	}
}

func (w *SqlWriter) WriteFile(p string, fs []string) {
	w.entry.file.path = p

	for i, f := range fs {
		w.entry.file.filters = append(w.entry.file.filters, sqlData{
			nr: i, value: f,
		})
	}
}

func (w *SqlWriter) WriteUser(u *user.User) {
	w.entry.user.login = u.Username
	w.entry.user.name = u.Name
}

func (w *SqlWriter) WriteTime(t, f time.Time) {
	w.entry.created = t.UTC()
	w.entry.file.modified = f.UTC()
}

func (w *SqlWriter) WriteHash(b []byte) {
	w.entry.file.hash = fmt.Sprintf("%x", b)
}

func (w *SqlWriter) WriteLines(ns []int, ss []string) {
	for i := range ss {
		w.entry.file.lines = append(w.entry.file.lines, sqlData{
			nr: ns[i], value: ss[i],
		})
	}
}

func (w *SqlWriter) insert(q string, v ...any) int64 {
	f := "?"

	for range len(v) - 1 {
		f += ", ?"
	}

	res, err := w.db.Exec(fmt.Sprintf("INSERT INTO %s VALUES (%s);", q, f), v...)

	if err != nil {
		sys.Error(err)
		return -1
	}

	id, err := res.LastInsertId()

	if err != nil {
		sys.Error(err)
		return -1
	}

	return id
}
