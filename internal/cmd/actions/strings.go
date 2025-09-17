package actions

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cuhsat/fox/internal"
	"github.com/cuhsat/fox/internal/app/ui"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/buffer"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
)

var StringsUsage = fox.Ascii + `
Display ASCII and Unicode strings.

Usage:
  fox strings [FLAG ...] PATH ...

Alias:
  carve, st

Positional arguments:
  Path(s) to open

Global:
  -p, --print              print directly to console
      --no-file            don't print filenames
      --no-line            don't print line numbers

Strings:
  -i, --ioc                detect built-in IoCs
  -e, --regexp=PATTERN     search for pattern
  -n, --min=NUMBER         minimum length (default: 3)
  -m, --max=NUMBER         maximum length (default: Unlimited)
  -a, --ascii              only carve ASCII strings

    Built-in IoC patterns:
      UUID, IPv4, IPv6, MAC, URL, Mail

Example:
  $ fox strings -in=8 malware.exe

Type "fox help" for more help...
`

var Strings = &cobra.Command{
	Use:     "strings",
	Short:   "display ASCII and Unicode strings",
	Long:    "display ASCII and Unicode strings",
	Args:    cobra.ArbitraryArgs,
	Aliases: []string{"carve", "st"},
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true

		re, err := cmd.Flags().GetString("regex")

		if len(re) > 0 {
			flg.Strings.Re, err = regexp.Compile(re)

			if err != nil {
				sys.Exit(fmt.Sprintf("could not compile regex: %s", err.Error()))
			}
		}

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
			os.Exit(2)
		} else if !flags.Get().Print {
			ui.Start(args, types.Strings)
		} else {
			flg := flags.Get()

			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(_ int, h *heap.Heap) {
				if h.Type != types.Stdin {
					if !flg.NoFile {
						fmt.Println(text.Title(h.String(), buffer.TermW))
					}

					for s := range h.Strings(
						flg.Strings.Min,
						flg.Strings.Max,
						flg.Strings.Ioc,
						flg.Strings.Re,
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
	Strings.Flags().BoolVarP(&flg.Strings.Ioc, "ioc", "i", false, "detect built-in IoCs")
	Strings.Flags().StringP("regex", "e", "", "search for specific pattern")
	Strings.Flags().IntVarP(&flg.Strings.Min, "min", "n", 3, "minimum length")
	Strings.Flags().IntVarP(&flg.Strings.Max, "max", "m", math.MaxInt, "maximum length")
	Strings.Flags().BoolVarP(&flg.Strings.Ascii, "ascii", "a", false, "only carve ASCII strings")
}
