package sys

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/pkg/text"
)

const (
	Dump = ".fox_dump"
)

var (
	vfs = make(map[string]File)
)

func Exit(v ...any) {
	_, _ = fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func Exec(cmds []string) File {
	f := TempFile("stdout")
	defer f.Close()

	for _, cmd := range cmds {
		args := text.SplitQuoted(cmd)

		if len(args) > 0 {
			var stdout, stderr bytes.Buffer

			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if cmd.Run() != nil {
				break
			}

			_, _ = f.WriteString(text.UnEscape(stdout.String()))
		}
	}

	return f
}

func Shell() {
	shell := os.Getenv("SHELL")

	if len(shell) == 0 {
		if runtime.GOOS == "windows" {
			shell = "CMD.EXE"
		} else {
			shell = "/bin/sh"
		}
	}

	fmt.Println(fox.Product, fox.Version)
	fmt.Println("Type 'exit' to return.")

	cmd := exec.Command(shell, "-l") // login shell
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()
}

func Stdin() File {
	if !IsPiped(os.Stdin) {
		Panic("invalid mode")
	}

	f := TempFile("stdin")

	go func(f File) {
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
				_ = f.Close()
				break

			default:
				Error(err)
			}
		}
	}(f)

	return f
}

func Stdout() File {
	return TempFile("stdout")
}

func Stderr() File {
	return TempFile("stderr")
}

func IsPiped(file File) bool {
	fi, err := file.Stat()

	if err != nil {
		Error(err)
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func Mapped(name string) File {
	return vfs[name]
}

func Exists(name string) bool {
	if Mapped(name) != nil {
		return true
	}

	_, err := os.Stat(name)

	return !errors.Is(err, os.ErrNotExist)
}

func Persist(name string) string {
	f := Mapped(name)

	if f == nil {
		return name
	}

	t, err := os.CreateTemp("", "fox-*")

	if err != nil {
		Panic(err)
	}

	_, err = f.WriteTo(t)

	if err != nil {
		Panic(err)
	}

	return t.Name()
}

func TempFile(name string) File {
	f := Mapped(name)

	if f != nil {
		return f
	}

	f = NewFileData(name)

	vfs[f.Name()] = f

	return f
}

func OpenFile(name string) File {
	f := Mapped(name)

	if f != nil {
		return f
	}

	f, err := os.OpenFile(name, os.O_RDONLY, 0400)

	if err != nil {
		Panic(err)
	}

	return f
}

func DumpStr(data string) File {
	f := TempFile("dump")
	defer f.Close()

	_, err := f.WriteString(data)

	if err != nil {
		Panic(err)
	}

	return f
}

func DumpErr(err any, stack any) {
	s := fmt.Sprintf("%+v\n\n%s", err, stack)

	err = os.WriteFile(Dump, []byte(s), 0600)

	if err != nil {
		Exit(err)
	}
}
