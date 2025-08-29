package sys

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/text"
)

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

			_, _ = f.WriteString(text.Unescape(stdout.String()))
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

func Map(file File) ([]byte, error) {
	b, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)

	if err != nil {
		return nil, err
	}

	return b, nil
}
