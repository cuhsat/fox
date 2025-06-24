package sys

import (
	"fmt"
	"log"
	"time"
)

var (
	Log *logger // global logger
)

type logger struct {
	file File // log file handle
}

func Setup() {
	Log = &logger{file: Stderr()}
	log.SetFlags(0)
	log.SetOutput(Log)
}

func (l logger) Name() string {
	return l.file.Name()
}

func (l logger) Write(b []byte) (int, error) {
	ts := time.Now().UTC().Format(time.RFC3339)

	return fmt.Fprintf(l.file, "[%s] %s", ts, string(b))
}

func (l logger) Close() {
	_ = l.file.Close()
}

func Error(v ...any) {
	log.Println(v...)
}

func Panic(v ...any) {
	log.Panic(v...)
}
