package sys

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	Prefix = "fox:"
)

var (
	Log *logger // global logger
)

type logger struct {
	f File // log file handle
}

func Setup() {
	Log = &logger{f: Stderr()}
	log.SetFlags(0)
	log.SetOutput(Log)
}

func (l logger) Name() string {
	return l.f.Name()
}

func (l logger) Close() {
	_ = l.f.Close()
}

func (l logger) Write(b []byte) (int, error) {
	_, _ = fmt.Fprint(os.Stderr, string(b))

	ts := time.Now().UTC().Format(time.RFC3339)

	return fmt.Fprintf(l.f, "[%s] %s", ts, string(b))
}

func Print(v ...any) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(Prefix+" %s", v...))
}

func Error(v ...any) {
	log.Println(fmt.Sprintf(Prefix+" %s", v...))
}

func Panic(v ...any) {
	log.Panic(fmt.Sprintf(Prefix+" %s", v...))
}
