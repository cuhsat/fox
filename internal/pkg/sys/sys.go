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

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/file"
)

const (
	Dump = ".fox_dump"
)

func Exit(v ...any) {
	Print(v...)
	os.Exit(1)
}

func Call(cmds []string) file.File {
	f := file.New("stdout")
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

	fmt.Println(app.Product, app.Version)
	fmt.Println("Type 'exit' to return.")

	cmd := exec.Command(shell, "-l") // login shell
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()
}

func Stdin() file.File {
	if !Piped(os.Stdin) {
		Panic("Device mode is invalid")
	}

	f := file.New("stdin")

	go func(f file.File) {
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

func Stdout() file.File {
	return file.New("stdout")
}

func Stderr() file.File {
	return file.New("stderr")
}

func Piped(file file.File) bool {
	fi, err := file.Stat()

	if err != nil {
		Error(err)
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func Open(name string) file.File {
	f := file.Open(name)

	if f != nil {
		return f // virtual file
	}

	f, err := os.OpenFile(name, os.O_RDONLY, 0400)

	if err == nil {
		return f // physical file
	}

	Panic(err)
	return nil
}

func TempFile() file.File {
	tmp, err := os.CreateTemp("", "fox-*")

	if err != nil {
		Panic(err)
	}

	return tmp
}

func TempDir() string {
	tmp, err := os.MkdirTemp("", "fox-*")

	if err != nil {
		Panic(err)
	}

	return tmp
}

func Exists(name string) bool {
	if file.Open(name) != nil {
		return true
	}

	_, err := os.Stat(name)

	return !errors.Is(err, os.ErrNotExist)
}

func Persist(name string) string {
	f := file.Open(name)

	if f == nil {
		return name // already persistent
	}

	t := TempFile()

	_, err := f.WriteTo(t)

	if err != nil {
		Panic(err)
	}

	return t.Name()
}

func Trace(err any, stack any) {
	if !flags.Get().Opt.Readonly {
		return // prevent dump
	}

	s := fmt.Sprintf("%+v\n\n%s", err, stack)

	err = os.WriteFile(Dump, []byte(s), 0600)

	if err != nil {
		Exit(err)
	}
}
