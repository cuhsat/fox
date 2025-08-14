// The Swiss Army Knife for examining text files.
//
// Usage of this code is governed by the GPL-3.0 License.
// Please see the LICENSE.md file for further information.
package main

import (
	"os"
	"runtime/debug"

	"github.com/hiforensics/fox/cmd"
	"github.com/hiforensics/fox/internal/pkg/sys"
)

// Starts the Forensic Examiner
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
