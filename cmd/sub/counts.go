package sub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

var Counts = &cobra.Command{
	Use:   "counts",
	Short: "display line and byte counts",
	Long:  "display line and byte counts",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		// force
		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true

		// invoke UI
		if !flg.Print {
			flg.UI.Invoke = types.Counts
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !flags.Get().Print {
			ui.Start(args)
		} else {
			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(h *heap.Heap) {
				fmt.Printf("%8dL %8dB  %s\n",
					h.Count(),
					h.Len(),
					h.String(),
				)
			})
		}
	},
}
