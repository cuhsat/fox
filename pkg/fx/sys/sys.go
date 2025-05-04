package sys

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/cuhsat/fx/pkg/fx"
)

const (
	FileDump = ".dump"
)

const (
	cmdShell = "CMD.EXE"
	pshShell = "/bin/sh"
)

type uncolor struct {
	File *os.File
	ansi *regexp.Regexp
}

func Exit(v ...any) {
	fmt.Fprintln(os.Stderr, v...)

	os.Exit(1)
}

func Exec(s string) string {
	uc := Uncolor(TempFile("exec", ".txt"))

	args := strings.Split(s, " ")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = uc
	cmd.Stderr = uc
	cmd.Run()

	uc.File.Close()

	return uc.File.Name()
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

	fmt.Println(fx.Product, fx.Version)
	fmt.Println("Type 'exit' to return.")

	cmd := exec.Command(shl, "-l") // login shell

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func Stdin() string {
	if !IsPiped(os.Stdin) {
		Panic("invalid mode")
	}

	f := TempFile("stdin", ".txt")

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
	return TempFile("stdout", ".txt")
}

func Stderr() *os.File {
	return TempFile("stderr", ".txt")
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

func TempFile(name, ext string) *os.File {
	f, err := os.CreateTemp("", fmt.Sprintf("fx-%s-*%s", name, ext))

	if err != nil {
		Panic(err)
	}

	return f
}

func DumpErr(err any, stack any) {
	s := fmt.Sprintf("%+v\n\n%s", err, stack)

	err = os.WriteFile(FileDump, []byte(s), 0600)

	if err != nil {
		Exit(err)
	}
}

func Uncolor(f *os.File) *uncolor {
	// remove 7-bit C1 ANSI sequences
	r := regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)

	return &uncolor{
		File: f,
		ansi: r,
	}
}

func (uc *uncolor) Write(p []byte) (n int, err error) {
	return uc.File.Write(uc.ansi.ReplaceAll(p, []byte("")))
}
