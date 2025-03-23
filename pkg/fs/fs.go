package fs

import (
    "fmt"
    "io"
    "os"
)

const (
    EX_OK      = 0
    EX_ERROR   = 1
    EX_USAGE   = 2
    EX_DATAERR = 3
    EX_NOINPUT = 4
)

const (
    MODE_FILE = 0644
)

const (
    FLAG_LOG = os.O_APPEND | os.O_CREATE | os.O_RDWR
    FLAG_FILE = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
)

func Error(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
}

func Panic(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
    os.Exit(EX_ERROR)
}

func Print(a ...any) {
    fmt.Fprintln(os.Stdout, a...)
    os.Exit(EX_OK)
}

func Usage(u string) {
    fmt.Fprintln(os.Stdout, "Usage:", u)
    os.Exit(EX_USAGE)
}

func Stdin(path string) {
    fi, err := os.Stdin.Stat()

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(EX_NOINPUT)
    }

    if (fi.Mode() & os.ModeCharDevice) != 0 {
        Panic("invalid mode")
    }

    b, err := io.ReadAll(os.Stdin)

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(EX_DATAERR)
    }

    err = os.WriteFile(path, b, MODE_FILE)

    if err != nil {
        Panic(err)
    }
}
