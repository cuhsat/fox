package sys

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hiforensics/fox/internal/pkg/types/file"
)

var (
	Log *logger // global logger
)

type logger struct {
	f file.File // log file handle
}

func Setup() {
	Log = &logger{f: Stderr()}
	log.SetFlags(0)
	log.SetOutput(Log)
}

func (l logger) Name() string {
	return l.f.Name()
}

func (l logger) Write(b []byte) (int, error) {
	ts := time.Now().UTC().Format(time.RFC3339)

	return fmt.Fprintf(l.f, "[%s] %s", ts, string(b))
}

func (l logger) Close() {
	_ = l.f.Close()
}

func Print(v ...any) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("fox: %s", v...))
}

func Error(v ...any) {
	log.Println(fmt.Sprintf("fox: %s", v...))
}

func Panic(v ...any) {
	log.Panic(fmt.Sprintf("fox: %s", v...))
}
