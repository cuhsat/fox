package fs

import (
    "fmt"
    "io"
    "os"
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

func Stdin() string {
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

    f, err := os.CreateTemp("", "cu-stdin-")

    if err != nil {
        Panic(err)
    }

    defer f.Close()

    _, err = f.Write(b)

    if err != nil {
        Panic(err)
    }

    return f.Name()
}
