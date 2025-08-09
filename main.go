package main

import (
	"os"
	"runtime/debug"

	"github.com/hiforensics/fox/cmd"
	"github.com/hiforensics/fox/internal/pkg/sys"
)

func main() {
	sys.Setup()

	defer func() {
		if err := recover(); err != nil {
			sys.Trace(err, debug.Stack())
			sys.Print(err)
		}

		sys.Log.Close()
	}()

	_ = os.Remove(sys.Dump)

	_ = cmd.Execute()
}
