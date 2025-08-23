// The Swiss Army Knife for examining text files.
package main

import (
	"os"
	"runtime/debug"

	"github.com/cuhsat/fox/internal/cmd"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

// Start fox.
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
