package sub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

var HashUsage = `
Display file hash sums.

Usage:
  fox hash [FLAG ...] PATH ...

Positional arguments:
  Path(s) to open

Global:
  -p, --print              print directly to console

Hash:
  -t, --type=ALGORITHM     hash algorithm (default: SHA256)

    Cryptographic hash algorithms:
      MD5, SHA1, SHA256, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

    Fuzzy hash algorithms:
      SDHASH, SSDEEP, TLSH

Example:
  $ fox hash -t=SHA3 artifacts.zip

Type "fox help" for more help...
`

var Hash = &cobra.Command{
	Use:   "hash",
	Short: "display file hash sums",
	Long:  "display file hash sums",
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !flags.Get().Print {
			ui.Start(args, types.Hash)
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

	Hash.SetHelpTemplate(fox.Fox + HashUsage)
	Hash.Flags().BoolVarP(&flg.Print, "print", "p", false, "print directly to console")
	Hash.Flags().VarP(&flg.Hash.Algo, "type", "t", "hash algorithm")
}
