package bag

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"time"

	_ "modernc.org/sqlite"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

const (
	schema = `
	CREATE TABLE IF NOT EXISTS album (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL
	)`
)

type SqlWriter struct {
	db *sql.DB // sql database
	tx *sql.Tx // sql transaction
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
		_, err = w.db.ExecContext(
			context.Background(),
			fmt.Sprintf("-- %s\n%s", t, schema),
		)

		if err != nil {
			sys.Panic(err)
		}
	}
}

func (w *SqlWriter) Start() {
	var err error

	w.tx, err = w.db.Begin()

	if err != nil {
		sys.Error(err)
	}
}

func (w *SqlWriter) Finalize() {
	err := w.tx.Commit()

	if err != nil {
		sys.Error(err)
	}
}

func (w *SqlWriter) WriteFile(p string, fs []string) {
	_, err := w.db.ExecContext(
		context.Background(),
		`INSERT INTO files (path) VALUES (?);`,
		p,
	)

	if err != nil {
		sys.Error(err)
	}

	// for _, f := range fs {
	// 	sb.WriteString(fmt.Sprintf(" > %s", f))
	// }
}

func (w *SqlWriter) WriteUser(u *user.User) {
	_, err := w.db.ExecContext(
		context.Background(),
		`INSERT INTO users (login, name) VALUES (?);`,
		u.Username,
		u.Name,
	)

	if err != nil {
		sys.Error(err)
	}
}

func (w *SqlWriter) WriteTime(t, f time.Time) {
	_, err := w.db.ExecContext(
		context.Background(),
		`INSERT INTO times (bagged, modified) VALUES (?);`,
		t.UTC(),
		f.UTC(),
	)

	if err != nil {
		sys.Error(err)
	}
}

func (w *SqlWriter) WriteHash(b []byte) {
	// fmt.Sprintf("%x\n", b)
}

func (w *SqlWriter) WriteLines(ns []int, ss []string) {
	// for i := 0; i < len(ss); i++ {
	// 	writeln(w.file, fmt.Sprintf("%08d  %v", ns[i], ss[i]))
	// }
}
