package sys

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

const (
    O_HISTORY = os.O_APPEND | os.O_CREATE | os.O_RDWR
    O_EVIDENCE = os.O_APPEND | os.O_CREATE | os.O_WRONLY
    O_OVERRIDE = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
)

func Debug(a ...any) {
    fmt.Fprintln(os.Stdout, a...)
}

func Print(a ...any) {
    fmt.Fprintln(os.Stdout, a...)
}

func Error(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
}

func Fatal(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
    os.Exit(1)
}

func Usage(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
    os.Exit(2)
}

func Stdin() string {
    if !IsPiped(os.Stdin) {
        Fatal("invalid mode")        
    }

    f := TempFile("stdin", "txt")

    go func(f *os.File) {
        r := bufio.NewReader(os.Stdin)

        for {
            s, err := r.ReadString('\n')

            switch err {
            case nil:
                _, err = f.WriteString(s)

                if err != nil {
                    Fatal(err)
                }

            case io.EOF:
                f.Close()
                break

            default:
                Fatal(err)
            }
        }
    }(f)

    return f.Name()
}

func Stdout() *os.File {
    return TempFile("stdout", "txt")
}

func Stderr() *os.File {
    return TempFile("stderr", "txt")
}

func IsPiped(f *os.File) bool {
    fi, err := f.Stat()

    if err != nil {
        Fatal(err)
    }

    return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func TempFile(n, e string) *os.File {
    f, err := os.CreateTemp("", fmt.Sprintf("fx-%s-*%s", n, e))

    if err != nil {
        Fatal(err)
    }

    return f
}
