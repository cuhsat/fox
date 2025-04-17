package sys

import (
    "bufio"
    "fmt"
    "io"
    "os"
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
    if !IsPipe(os.Stdin) {
        Fatal("invalid mode")        
    }

    f := Temp("stdin", "txt")

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
    return Temp("stdout", "txt")
}

func Stderr() *os.File {
    return Temp("stderr", "txt")
}

func IsPipe(f *os.File) bool {
    fi, err := f.Stat()

    if err != nil {
        Fatal(err)
    }

    return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func Temp(n, e string) *os.File {
    f, err := os.CreateTemp("", fmt.Sprintf("fx-%s-*%s", n, e))

    if err != nil {
        Fatal(err)
    }

    return f
}

func Open(p string) *os.File {
    f, err := os.OpenFile(p, os.O_RDONLY, 0400)

    if err != nil {
        Fatal(err)
    }

    return f
}
