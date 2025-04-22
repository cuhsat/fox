package fx

import (
    "fmt"
    "log"
    "os"
    "time"
)

var (
    Log *logger // default logger
)

type logger struct {
    Name string   // log file name
    file *os.File // log file handle
}

func SetupLogger() {
    f := Stderr()

    Log = &logger{
        Name: f.Name(),
        file: f,
    }

    log.SetFlags(0)
    log.SetOutput(Log)
}

func (l logger) Write(b []byte) (int, error) {
    ts := time.Now().UTC().Format(time.RFC3339)

    return fmt.Fprintf(l.file, "[%s] %s", ts, string(b))
}

func (l logger) Close() {
    l.file.Close()

    os.Remove(l.Name)
}

func Error(v ...any) {
    log.Println(v...)
}

func Panic(v ...any) {
    log.Panic(v...)
}
