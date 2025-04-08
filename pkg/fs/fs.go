package fs

import (
    "bufio"
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
        Panic(err)
    }

    if (fi.Mode() & os.ModeCharDevice) != 0 {
        Panic("invalid mode")
    }

    f, err := os.CreateTemp("", "cu-stdin-")

    if err != nil {
        Panic(err)
    }

    go func(f *os.File) {
        r := bufio.NewReader(os.Stdin)

        for {
            s, err := r.ReadString('\n')

            switch err {
            case nil:
                _, err = f.WriteString(s)

                if err != nil {
                    Panic(err)
                }

            case io.EOF:
                f.Close()
                break

            default:
                Panic(err)
            }
        }
    }(f)

    return f.Name()
}
