package fs

import (
    "fmt"
    "io"
    "os"
    "os/exec"

    "github.com/mattn/go-shellwords"
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

func Print(a ...any) {
    fmt.Fprintln(os.Stdout, a...)
    os.Exit(0)
}

func Usage(u string) {
    fmt.Fprintln(os.Stdout, "Usage:", u)
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

func System(cmd string) error {
    arg, err := shellwords.Parse(cmd)

    if err != nil {
        return err
    }

    _, err = exec.Command(arg[0], arg[1:]...).Output()

    return err
}
