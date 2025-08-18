package actions

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/app"
	"github.com/hiforensics/fox/internal/app/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

var CountsUsage string = app.Ascii + `
Display line and byte counts.

Usage:
  fox counts [FLAG ...] PATH ...

Positional arguments:
  Path(s) to open

Global:
  -p, --print              print directly to console

Example:
  $ fox counts ./**/*.txt

Type "fox help" for more help...
`

var Counts = &cobra.Command{
	Use:   "counts",
	Short: "display line and byte counts",
	Long:  "display line and byte counts",
	Args:  cobra.ArbitraryArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print(CountsUsage)
			os.Exit(0)
		} else if !flags.Get().Print {
			ui.Start(args, types.Counts)
		} else {
			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(h *heap.Heap) {
				fmt.Printf("%8dL %8dB  %s\n", h.Count(), h.Len(), h.String())
			})
		}
	},
}

func init() {
	flg := flags.Get()

	Counts.SetHelpTemplate(CountsUsage)
	Counts.Flags().BoolVarP(&flg.Print, "print", "p", false, "print directly to console")
}
