package bag

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"time"

	_ "embed"
	_ "modernc.org/sqlite"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

var (
	//go:embed bag.sql
	schema string
)

type sqlEvidence struct {
	user   sqlUser
	file   sqlFile
	bagged time.Time
}

type sqlUser struct {
	login string
	name  string
}

type sqlFile struct {
	path     string
	hash     string
	modified time.Time
	filters  []sqlFilter
	lines    []sqlLine
}

type sqlFilter struct {
	nr    int
	value string
}

type sqlLine struct {
	nr    int
	value string
}

type SqlWriter struct {
	db    *sql.DB      // sql database
	tx    *sql.Tx      // sql transaction
	entry *sqlEvidence // current entry
}

func NewSqlWriter() *SqlWriter {
	return &SqlWriter{
		db: nil,
		tx: nil,
	}
}

func (w *SqlWriter) Init(f *os.File, n bool, t string) {
	var err error

	f.Close()

	w.db, err = sql.Open("sqlite", fmt.Sprintf("file:%s", f.Name()))

	if err != nil {
		sys.Error(err)
	}

	if n {
		_, err = w.db.Exec(schema)

		if err != nil {
			sys.Panic(err)
		}
	}
}

func (w *SqlWriter) Start() {
	w.entry = new(sqlEvidence)
}

func (w *SqlWriter) Finalize() {
	var err error

	w.tx, err = w.db.Begin()

	if err != nil {
		sys.Error(err)
	}

	res, err := w.db.Exec(
		`INSERT INTO users (login, name) VALUES (?, ?);`,
		w.entry.user.login,
		w.entry.user.name,
	)

	if err != nil {
		sys.Error(err)
	}

	user_id, err := res.LastInsertId()

	if err != nil {
		sys.Error(err)
	}

	res, err = w.db.Exec(
		`INSERT INTO files (path, hash, modified) VALUES (?, ?, ?);`,
		w.entry.file.path,
		w.entry.file.hash,
		w.entry.file.modified,
	)

	if err != nil {
		sys.Error(err)
	}

	file_id, err := res.LastInsertId()

	if err != nil {
		sys.Error(err)
	}

	for _, f := range w.entry.file.filters {
		if _, err := w.db.Exec(
			`INSERT INTO filters (file_id, nr, value) VALUES (?, ?, ?);`,
			file_id,
			f.nr,
			f.value,
		); err != nil {
			sys.Error(err)
		}
	}

	for _, l := range w.entry.file.lines {
		if _, err := w.db.Exec(
			`INSERT INTO lines (file_id, nr, value) VALUES (?, ?, ?);`,
			file_id,
			l.nr,
			l.value,
		); err != nil {
			sys.Error(err)
		}
	}

	if _, err := w.db.Exec(
		`INSERT INTO evidence (user_id, file_id, bagged) VALUES (?, ?, ?);`,
		user_id,
		file_id,
		w.entry.bagged,
	); err != nil {
		sys.Error(err)
	}

	err = w.tx.Commit()

	if err != nil {
		sys.Error(err)
	}
}

func (w *SqlWriter) WriteFile(p string, fs []string) {
	w.entry.file = sqlFile{
		path: p,
	}

	for i, f := range fs {
		w.entry.file.filters = append(w.entry.file.filters, sqlFilter{
			nr: i, value: f,
		})
	}
}

func (w *SqlWriter) WriteUser(u *user.User) {
	w.entry.user = sqlUser{
		login: u.Username, name: u.Name,
	}
}

func (w *SqlWriter) WriteTime(t, f time.Time) {
	w.entry.bagged = t.UTC()
	w.entry.file.modified = f.UTC()
}

func (w *SqlWriter) WriteHash(b []byte) {
	w.entry.file.hash = fmt.Sprintf("%x", b)
}

func (w *SqlWriter) WriteLines(ns []int, ss []string) {
	for i := range ss {
		w.entry.file.lines = append(w.entry.file.lines, sqlLine{
			nr: ns[i], value: ss[i],
		})
	}
}
