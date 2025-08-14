package sub

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

// Fox deflate usage
var DeflateUsage = fox.Fox + `
Deflate compressed files.

Usage:
  fox deflate [FLAG ...] PATH...

Positional arguments:
  Path(s) to open

Global:
      --no-file            don't print filenames

Deflate:
  -d, --dir[=PATH]         deflate into directory (default: .)
      --pass=PASSWORD      decrypt with password (RAR and ZIP only)

Example:
  $ fox deflate --pass=infected ioc.rar

Type "fox help" for more help...
`

// Deflates compressed files
var Deflate = &cobra.Command{
	Use:   "deflate",
	Short: "deflate compressed files",
	Long:  "deflate compressed files",
	Args:  cobra.ArbitraryArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print(DeflateUsage)
			os.Exit(0)
		}

		flg := flags.Get()

		hs := heapset.New(args)
		defer hs.ThrowAway()

		hs.Each(func(h *heap.Heap) {
			root := flg.Deflate.Path

			if root == "." {
				name := filepath.Base(h.Base)
				root = name[0 : len(name)-len(filepath.Ext(name))]
			}

			// convert to relative path
			path := h.Title

			if h.Type == types.Deflate {
				path = path[len(h.Base)+1:]
			} else {
				path = filepath.Base(path)
			}

			// create (sub)folders
			if sub := filepath.Dir(path); len(sub) > 0 {
				sub = filepath.Join(root, sub)

				err := os.MkdirAll(sub, 0700)

				if err != nil {
					sys.Exit(err)
				}
			}

			path = filepath.Join(root, path)

			if !flg.NoFile {
				fmt.Printf("Deflate %s\n", path)
			}

			err := os.WriteFile(path, *h.MMap(), 0600)

			if err != nil {
				sys.Exit(err)
			}
		})

		fmt.Printf("%d file(s) written\n", hs.Len())
	},
}

func init() {
	flg := flags.Get()

	Deflate.SetHelpTemplate(DeflateUsage)
	Deflate.Flags().StringVarP(&flg.Deflate.Path, "dir", "d", "", "deflate into directory")
	Deflate.Flags().Lookup("dir").NoOptDefVal = "."
	Deflate.MarkFlagDirname("dir")
}
