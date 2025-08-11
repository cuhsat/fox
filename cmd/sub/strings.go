package sub

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/text"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/buffer"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

var StringsUsage = `
Display ASCII and Unicode strings.

Usage:
  fox strings [FLAG...] [PATH...]

Positional arguments:
  Path(s) to open

Global:
  -p, --print              print directly to console
      --no-file            don't print filenames
      --no-line            don't print line numbers

Strings:
  -n, --min=NUMBER         minimum length (default: 3)
  -m, --max=NUMBER         maximum length (default: Unlimited)

Example:
  $ fox strings -n=8 malware.exe

Type "fox help" for more help...
`

var Strings = &cobra.Command{
	Use:   "strings",
	Short: "display ASCII and Unicode strings",
	Long:  "display ASCII and Unicode strings",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		// force
		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !flags.Get().Print {
			ui.Start(args, types.Strings)
		} else {
			flg := flags.Get()

			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(h *heap.Heap) {
				if h.Type != types.Stdin {
					if !flg.NoFile {
						fmt.Println(text.Title(h.String(), buffer.TermW))
					}

					for s := range h.Strings(
						flg.Strings.Min,
						flg.Strings.Max,
					) {
						if !flg.NoLine {
							fmt.Printf("%08x  %s\n", s.Off, strings.TrimSpace(s.Str))
						} else {
							fmt.Println(strings.TrimSpace(s.Str))
						}
					}
				}
			})
		}
	},
}

func init() {
	flg := flags.Get()

	Strings.SetHelpTemplate(fox.Fox + StringsUsage)
	Strings.Flags().BoolVarP(&flg.Print, "print", "p", false, "print directly to console")
	Strings.Flags().BoolVarP(&flg.NoFile, "no-file", "", false, "don't print filenames")
	Strings.Flags().BoolVarP(&flg.NoLine, "no-line", "", false, "don't print line numbers")
	Strings.Flags().IntVarP(&flg.Strings.Min, "min", "n", 3, "minimum length")
	Strings.Flags().IntVarP(&flg.Strings.Max, "max", "m", 0, "maximum length")
}
