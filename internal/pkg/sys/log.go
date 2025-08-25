package sys

import (
	"fmt"
	"log"
	"os"
	"time"
)

const Prefix = "fox:"

var Log *logger

func Setup() {
	Log = &logger{f: Stderr()}
	log.SetFlags(0)
	log.SetOutput(Log)
}

type logger struct {
	f File // log file handle
}

func (l logger) Name() string {
	return l.f.Name()
}

func (l logger) Write(b []byte) (int, error) {
	_, _ = fmt.Fprint(os.Stderr, string(b))

	ts := time.Now().UTC().Format(time.RFC3339)

	return fmt.Fprintf(l.f, "[%s] %s", ts, string(b))
}

func Trace(v any, stack any) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(Prefix+" %+v\n\n%s", v, stack))
}

func Debug(v ...any) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(Prefix+" %#v", v...))
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
