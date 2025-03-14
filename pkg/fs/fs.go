package fs

import (
    "fmt"
    "os"
)

const (
    EX_OK    = 0
    EX_ERROR = 1
    EX_USAGE = 2
)

const (
    MODE_FILE = 0644
)

func Error(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
}

func Panic(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
    os.Exit(EX_ERROR)
}

func Usage(u string) {
    fmt.Fprintln(os.Stdout, "Usage:", u)
    os.Exit(EX_USAGE)
}
