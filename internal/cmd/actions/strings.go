package actions

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/app"
	"github.com/hiforensics/fox/internal/app/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/text"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/buffer"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

var StringsUsage = app.Ascii + `
Display ASCII and Unicode strings.

Usage:
  fox strings [FLAG ...] PATH ...

Positional arguments:
  Path(s) to open

Global:
  -p, --print              print directly to console
      --no-file            don't print filenames
      --no-line            don't print line numbers

Strings:
  -a, --ascii              only ASCII strings
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
	Args:  cobra.ArbitraryArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true

		if flg.Strings.Min <= 0 {
			sys.Exit("min must be greater than 0")
		}

		if flg.Strings.Max <= 0 {
			sys.Exit("min must be greater than 0")
		}

		if flg.Strings.Min > flg.Strings.Max {
			sys.Exit("max must be greater than min")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print(StringsUsage)
			os.Exit(0)
		} else if !flags.Get().Print {
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

	Strings.SetHelpTemplate(StringsUsage)
	Strings.Flags().BoolVarP(&flg.Print, "print", "p", false, "print directly to console")
	Strings.Flags().BoolVarP(&flg.NoFile, "no-file", "", false, "don't print filenames")
	Strings.Flags().BoolVarP(&flg.NoLine, "no-line", "", false, "don't print line numbers")
	Strings.Flags().IntVarP(&flg.Strings.Min, "min", "n", 3, "minimum length")
	Strings.Flags().IntVarP(&flg.Strings.Max, "max", "m", math.MaxInt, "maximum length")
	Strings.Flags().BoolVarP(&flg.Strings.Ascii, "ascii", "a", false, "only ASCII strings")
}
