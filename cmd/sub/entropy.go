package sub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

var EntropyUsage string = `
Display file entropy.

Usage:
  fox entropy [FLAG ...] PATH ...

Positional arguments:
  Path(s) to open

Global:
  -p, --print              print directly to console

Example:
  $ fox entropy ./**/*

Type "fox help" for more help...
`

var Entropy = &cobra.Command{
	Use:   "entropy",
	Short: "display file entropy",
	Long:  "display file entropy",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		// force
		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !flags.Get().Print {
			ui.Start(args, types.Entropy)
		} else {
			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(h *heap.Heap) {
				fmt.Printf("%.10f  %s\n", h.Entropy(), h.String())
			})
		}
	},
}

func init() {
	flg := flags.Get()

	Entropy.SetHelpTemplate(fox.Fox + EntropyUsage)
	Entropy.Flags().BoolVarP(&flg.Print, "print", "p", false, "print directly to console")
}
