package sys

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	mem "github.com/cuhsat/memfile"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/text"
)

var fs = make(map[string]*mem.File)

type File = mem.Fileable

func Exit(v ...any) {
	Print(v...)
	os.Exit(1)
}

func Call(cmds []string) File {
	f := Create("stdout")
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

func Stdin() File {
	if !Piped(os.Stdin) {
		Panic("Device mode is invalid")
	}

	f := Create("stdin")

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
	return Create("stdout")
}

func Stderr() File {
	return Create("stderr")
}

func Piped(file File) bool {
	fi, err := file.Stat()

	if err != nil {
		Error(err)
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func Open(name string) File {
	vtl, ok := fs[name]

	if ok {
		vtl.Seek(0, io.SeekStart)
		return vtl // virtual file
	}

	reg, err := os.OpenFile(name, os.O_RDONLY, 0400)

	if err == nil {
		return reg // regular file
	}

	Panic(err)
	return nil
}

func Create(name string) File {
	uniq := fmt.Sprintf("fox://%d/%s", rand.Uint64(), name)
	file := mem.New(uniq)

	fs[uniq] = file

	return file
}

func Exists(name string) bool {
	_, ok := fs[name]

	if ok {
		return true
	}

	_, err := os.Stat(name)

	return !errors.Is(err, os.ErrNotExist)
}

func Persist(name string) string {
	f, ok := fs[name]

	if !ok {
		return name // already persistent
	}

	t, err := os.CreateTemp(Cache(), "fox-*")

	if err != nil {
		Panic(err)
	}

	_, err = f.WriteTo(t)

	if err != nil {
		Panic(err)
	}

	return t.Name()
}

func Config(name string) string {
	dir, err := os.UserHomeDir()

	if err != nil {
		Panic(err)
	}

	return filepath.Join(dir, ".config", "fox", name)
}

func Cache() string {
	dir, err := os.UserHomeDir()

	if err != nil {
		Panic(err)
	}

	tmp := filepath.Join(dir, ".cache", "fox")

	err = os.MkdirAll(tmp, 0700)

	if err != nil {
		Panic(err)
	}

	return tmp
}
