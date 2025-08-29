package actions

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/app/ui"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
)

var EntropyUsage = app.Ascii + `
Display file entropy.

Usage:
  fox entropy [FLAG ...] PATH ...

Alias:
	en

Positional arguments:
  Path(s) to open

Global:
  -p, --print              print directly to console

Entropy:
  -n, --min[=DECIMAL]      minimum score (default: 0.8)
  -m, --max[=DECIMAL]      maximum score (default: 0.8)

Example:
  $ fox entropy -n ./**/*

Type "fox help" for more help...
`

var Entropy = &cobra.Command{
	Use:     "entropy",
	Short:   "display file entropy",
	Long:    "display file entropy",
	Aliases: []string{"en"},
	Args:    cobra.ArbitraryArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true

		if flg.Entropy.Min < 0 {
			sys.Exit("min must be 0 or greater")
		}

		if flg.Entropy.Max > 1 {
			sys.Exit("max must be 1 or lesser")
		}

		if flg.Entropy.Min > flg.Entropy.Max {
			sys.Exit("max must be greater than min")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print(EntropyUsage)
			os.Exit(0)
		} else if !flags.Get().Print {
			ui.Start(args, types.Entropy)
		} else {
			flg := flags.Get()

			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(_ int, h *heap.Heap) {
				if v := h.Entropy(
					flg.Entropy.Min,
					flg.Entropy.Max,
				); v != -1 {
					fmt.Printf("%.10f  %s\n", v, h.String())
				}
			})
		}
	},
}

func init() {
	flg := flags.Get()

	Entropy.SetHelpTemplate(EntropyUsage)
	Entropy.Flags().BoolVarP(&flg.Print, "print", "p", false, "print directly to console")
	Entropy.Flags().Float64VarP(&flg.Entropy.Min, "min", "n", 0.0, "minimum score")
	Entropy.Flags().Float64VarP(&flg.Entropy.Max, "max", "m", 1.0, "maximum score")
	Entropy.Flags().Lookup("min").NoOptDefVal = "0.8"
	Entropy.Flags().Lookup("max").NoOptDefVal = "0.8"
}
