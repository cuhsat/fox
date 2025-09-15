package sys

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/rainu/go-command-chain"
	"github.com/spf13/afero"

	"github.com/cuhsat/fox/internal"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
)

type File = afero.File

func Exit(v ...any) {
	Print(v...)
	os.Exit(1)
}

func Exec(cmds []string) File {
	f := fs.Create("/fox/exec")
	defer f.Close()

	for _, cmd := range cmds {
		err := cmdchain.Builder().JoinShellCmd(cmd).
			Finalize().WithOutput(f).Run()

		if err != nil {
			Error(err)
			break
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
	if !Piped(os.Stdin) {
		Panic("Device mode is invalid")
	}

	f := fs.Create("/fox/stdin")

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
	return fs.Create("/fox/stdout")
}

func Stderr() File {
	return fs.Create("/fox/stderr")
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
