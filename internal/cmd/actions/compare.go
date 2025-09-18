package actions

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cuhsat/fox/internal"
	"github.com/cuhsat/fox/internal/app/ui"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
)

var CompareUsage = fox.Ascii + `
Compare two files.

Usage:
  fox compare [FLAG ...] FILE1 FILE2

Alias:
  diff, cmp, ce

Positional arguments:
  Files to open

Global:
  -p, --print              print directly to console
      --no-file            don't print filenames
      --no-line            don't print line numbers

Compare:
  -g, --git                use the unified git diff format

Example:
  $ fox compare server.log mirror.log

Type "fox help" for more help...
`

var Compare = &cobra.Command{
	Use:     "compare",
	Short:   "compare two files",
	Long:    "compare two files",
	Args:    cobra.ArbitraryArgs,
	Aliases: []string{"diff", "cmp", "ce"},
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Print(CompareUsage)
			os.Exit(2)
		} else if !flags.Get().Print {
			ui.Start(args, types.Compare)
		} else {
			var a [2]*heap.Heap

			hs := heapset.New(args)
			defer hs.ThrowAway()

			hs.Each(func(i int, h *heap.Heap) {
				a[i] = h
			})

			fmt.Print(text.Diff(
				a[0].String(),
				a[1].String(),
				a[0].SMap().Lines(),
				a[1].SMap().Lines(),
				flags.Get().Compare.Git,
			))
		}
	},
}

func init() {
	flg := flags.Get()

	Compare.SetHelpTemplate(CompareUsage)
	Compare.Flags().BoolVarP(&flg.Print, "print", "p", false, "print directly to console")
	Compare.Flags().BoolVarP(&flg.NoFile, "no-file", "", false, "don't print filenames")
	Compare.Flags().BoolVarP(&flg.NoLine, "no-line", "", false, "don't print line numbers")
	Compare.Flags().BoolVarP(&flg.Compare.Git, "git", "g", false, "use the unified git diff format")
}
