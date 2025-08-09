package sub

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
)

var Deflate = &cobra.Command{
	Use:   "deflate",
	Short: "deflate compressed files",
	Long:  "deflate compressed files",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		// force
		flg.Opt.NoConvert = true
		flg.Opt.NoPlugins = true
	},
	Run: func(cmd *cobra.Command, args []string) {
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

			err := os.WriteFile(path, *h.Ensure().MMap(), 0600)

			if err != nil {
				sys.Exit(err)
			}
		})

		fmt.Printf("%d file(s) written\n", hs.Len())
	},
}

func init() {
	flg := flags.Get()

	Deflate.Flags().StringVarP(&flg.Deflate.Path, "dir", "d", "", "deflate into directory")
	Deflate.Flags().Lookup("dir").NoOptDefVal = "."
	Deflate.MarkFlagDirname("dir")
}
