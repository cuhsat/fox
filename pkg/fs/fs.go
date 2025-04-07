package fs

import (
    "fmt"
    "io"
    "os"
)

const (
    In  = ".cin"
    Out = ".cout"
)

const (
    Append   = os.O_CREATE | os.O_APPEND | os.O_RDWR
    Override = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
)

func Error(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
}

func Panic(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
    os.Exit(1)
}

func Usage(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
    os.Exit(2)
}

func Stdin(path string) {
    fi, err := os.Stdin.Stat()

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(4)
    }

    if (fi.Mode() & os.ModeCharDevice) != 0 {
        Panic("invalid mode")
    }

    b, err := io.ReadAll(os.Stdin)

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(3)
    }

    err = os.WriteFile(path, b, 0644)

    if err != nil {
        Panic(err)
    }
}
