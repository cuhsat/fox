package sub

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/text"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/buffer"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

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

		// invoke UI
		if !flg.Print {
			flg.UI.Invoke = types.Strings
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !flags.Get().Print {
			ui.Start(args)
		} else {
			flg := flags.Get()

			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(h *heap.Heap) {
				if h.Type != types.Stdin {
					if !flg.NoFile {
						fmt.Println(text.Title(h.String(), buffer.TermW))
					}

					for s := range h.Strings(flg.Strings.Min) {
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

	Strings.Flags().IntVarP(&flg.Strings.Min, "min", "m", 3, "minimum length")
}
