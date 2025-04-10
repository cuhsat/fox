package fs

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

const (
    O_HISTORY = os.O_APPEND | os.O_CREATE | os.O_RDWR
    O_OVERRIDE = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
)

func Print(a ...any) {
    fmt.Fprintln(os.Stdout, a...)
}

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

    f := TempFile("stdin")

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

func Stdout() *os.File {
    return TempFile("stdout")
}

func Stderr() *os.File {
    return TempFile("stderr")
}

func IsStdout() bool {
    fi, err := os.Stdout.Stat()

    if err != nil {
        Panic(err)
    }

    return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func TempFile(s string) *os.File {
    f, err := os.CreateTemp("", fmt.Sprintf("cu-%s-", s))

    if err != nil {
        Panic(err)
    }

    return f
}
