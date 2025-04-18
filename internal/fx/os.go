package fx

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "os"
)

var (
    Logfile *os.File
)

func Init() *os.File {
    f := Stderr()

    log.SetOutput(f)

    Logfile = f

    return f
}

func Error(a ...any) {
    log.Println(a...)
}

func Fatal(a ...any) {
    log.Fatal(a...)
}

func Stdin() string {
    if !IsPiped(os.Stdin) {
        Fatal("invalid mode")
    }

    f := Temp("stdin", ".txt")

    go func(f *os.File) {
        r := bufio.NewReader(os.Stdin)

        for {
            s, err := r.ReadString('\n')

            switch err {
            case nil:
                _, err = f.WriteString(s)

                if err != nil {
                    Error(err)
                }

            case io.EOF:
                f.Close()
                break

            default:
                Error(err)
            }
        }
    }(f)

    return f.Name()
}

func Stdout() *os.File {
    return Temp("stdout", ".txt")
}

func Stderr() *os.File {
    return Temp("stderr", ".txt")
}

func IsPiped(f *os.File) bool {
    fi, err := f.Stat()

    if err != nil {
        Error(err)

        return false
    }

    is := (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice

    return is
}

func Open(path string) *os.File {
    f, err := os.OpenFile(path, os.O_RDONLY, 0400)

    if err != nil {
        Fatal(err)
    }

    return f
}

func Temp(name, ext string) *os.File {
    f, err := os.CreateTemp("", fmt.Sprintf("fx-%s-*%s", name, ext))

    if err != nil {
        Fatal(err)
    }

    return f
}
