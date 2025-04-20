package fx

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
)

var (
    Logfile string
)

func Init() *os.File {
    f := Stderr()

    log.SetFlags(log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
    log.SetOutput(f)

    Logfile = f.Name()

    return f
}

func Exit() {
    if len(Logfile) > 0 {
        os.Remove(Logfile)
    }
}

func Debug(v ...any) {
    log.Println(v...)
}

func Error(v ...any) {
    log.Println(v...)
}

func Fatal(v ...any) {
    log.Fatal(v...)
}

func Panic(v ...any) {
    log.Panic(v...)
}

func Stdin() string {
    if !IsPiped(os.Stdin) {
        Panic("invalid mode")
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

func Exists(path string) bool {
    _, err := os.Stat(path)

   return !errors.Is(err, os.ErrNotExist)
}

func Open(path string) *os.File {
    f, err := os.OpenFile(path, os.O_RDONLY, 0400)

    if err != nil {
        Panic(err)
    }

    return f
}

func Temp(name, ext string) *os.File {
    f, err := os.CreateTemp("", fmt.Sprintf("fx-%s-*%s", name, ext))

    if err != nil {
        Panic(err)
    }

    return f
}

func Dump(err any, stack any) {
    s := fmt.Sprintf("%+v\n\n%s", err, stack)

    err = os.WriteFile(".dump", []byte(s), 0600)

    if err != nil {
        Fatal(err)
    }
}
