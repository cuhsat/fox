package sys

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/pkg/text"
)

const (
	FileDump = ".dump"
)

const (
	cmdShell = "CMD.EXE"
	pshShell = "/bin/sh"
)

func Exit(v ...any) {
	_, _ = fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func Exec(cmds []string) string {
	f := TempFile()
	defer f.Close()

	var stdout io.ReadCloser = nil
	var stderr io.ReadCloser = nil

	for _, cmd := range cmds {
		args := text.Split(cmd)

		if len(args) > 0 {
			cmd := exec.Command(args[0], args[1:]...)

			if stdout != nil && stderr != nil {
				cmd.Stdin = io.MultiReader(stderr, stdout)
			}

			stdout, _ = cmd.StdoutPipe()
			stderr, _ = cmd.StderrPipe()

			_ = cmd.Start()
			defer cmd.Wait()
		}
	}

	if stdout != nil {
		go io.Copy(f, stdout)
	}

	if stderr != nil {
		go io.Copy(f, stderr)
	}

	return f.Name()
}

func Shell() {
	shl := os.Getenv("SHELL")

	if len(shl) == 0 {
		if runtime.GOOS == "windows" {
			shl = cmdShell
		} else {
			shl = pshShell
		}
	}

	fmt.Println(fox.Product, fox.Version)
	fmt.Println("Type 'exit' to return.")

	cmd := exec.Command(shl, "-l") // login shell

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()
}

func Stdin() string {
	if !IsPiped(os.Stdin) {
		Panic("invalid mode")
	}

	f := TempFile()

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
				_ = f.Close()
				break

			default:
				Error(err)
			}
		}
	}(f)

	return f.Name()
}

func Stdout() *os.File {
	return TempFile()
}

func Stderr() *os.File {
	return TempFile()
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

func OpenFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_RDONLY, 0400)

	if err != nil {
		Panic(err)
	}

	return f
}

func TempFile() *os.File {
	f, err := os.CreateTemp("", "fox-*")

	if err != nil {
		Panic(err)
	}

	return f
}

func Extract(data string) string {
	f := TempFile()

	_, err := f.WriteString(data)

	if err != nil {
		Error(err)
	}

	_ = f.Close()

	return f.Name()
}

func DumpErr(err any, stack any) {
	s := fmt.Sprintf("%+v\n\n%s", err, stack)

	err = os.WriteFile(FileDump, []byte(s), 0600)

	if err != nil {
		Exit(err)
	}
}
