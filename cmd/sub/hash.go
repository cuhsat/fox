package sub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

var Hash = &cobra.Command{
	Use:   "hash",
	Short: "display hash sums",
	Long:  "display hash sums",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		// force
		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true

		// default
		if len(flg.Hash.Algo) == 0 {
			flg.Hash.Algo = types.SHA256
		}

		// invoke UI
		if !flg.Print {
			flg.UI.Invoke = types.Hash
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !flags.Get().Print {
			ui.Start(args)
		} else {
			algo := flags.Get().Hash.Algo.String()

			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(h *heap.Heap) {
				sum, err := h.HashSum(algo)

				if err != nil {
					sys.Exit(err)
				}

				switch algo {
				case types.SDHASH:
					fmt.Printf("%s  %s\n", sum, h.String())
				default:
					fmt.Printf("%x  %s\n", sum, h.String())
				}
			})
		}
	},
}

func init() {
	flg := flags.Get()

	Hash.Flags().VarP(&flg.Hash.Algo, "type", "t", "hash algorithm")
}
