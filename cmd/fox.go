package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/hiforensics/fox/cmd/sub"
	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/text"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/buffer"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
	"github.com/hiforensics/fox/internal/pkg/types/mode"
	"github.com/spf13/cobra"
)

var Fox = &cobra.Command{
	Use:     "fox",
	Short:   "examine text files",
	Long:    "examine text files",
	Args:    cobra.ArbitraryArgs,
	Version: fox.Version,
	PreRun: func(cmd *cobra.Command, args []string) {
		flg := flags.Get()

		// print if output is piped
		if sys.Piped(os.Stdout) {
			flg.Print = true
		}

		if flg.Filters.Context > 0 {
			flg.Filters.Before = flg.Filters.Context
			flg.Filters.After = flg.Filters.Context
		}

		if flg.Opt.Raw {
			flg.Opt.NoConvert = true
			flg.Opt.NoDeflate = true
			flg.Opt.NoPlugins = true
		}

		if flg.Bag.No {
			flg.Bag.Mode = flags.BagModeNone
		}

		if flg.Alias.Text {
			flg.Bag.Mode = flags.BagModeText
		}

		if flg.Alias.Json {
			flg.Bag.Mode = flags.BagModeJson
		}

		if flg.Alias.Jsonl {
			flg.Bag.Mode = flags.BagModeJsonl
		}

		if flg.Alias.Xml {
			flg.Bag.Mode = flags.BagModeXml
		}

		if flg.Alias.Sqlite {
			flg.Bag.Mode = flags.BagModeSqlite
		}

		// explicit set UI mode
		if flg.Hex {
			flg.UI.Mode = mode.Hex
		}

		// implicit set UI mode
		if len(flg.Filters.Patterns) > 0 {
			flg.UI.Mode = mode.Grep
		}

		if len(flg.UI.State) > 0 {
			re := regexp.MustCompile("[^-nwtNWT]+")

			flg.UI.State = re.ReplaceAllString(flg.UI.State, "")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !flags.Get().Print {
			ui.Start(args)
		} else {
			print(args)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if l := sys.Log.Consume(); len(l) > 0 {
			log.Print(l)
		}
	},
}

func Execute() error {
	return Fox.Execute()
}

func print(args []string) {
	flg := flags.Get()

	hs := heapset.New(args)
	defer hs.ThrowAway()

	hs.Each(func(h *heap.Heap) {
		if h.Type != types.Stdin {
			ctx := buffer.NewContext(h)

			if !flg.NoFile {
				fmt.Println(text.Title(h.String(), buffer.TermW))
			}

			if len(flg.Filters.Patterns) > 0 {
				if ctx.Heap.Len() == 0 {
					return // ignore empty files
				}

				flg := flags.Get()

				for l := range buffer.Text(ctx).Lines {
					if l.Nr == "--" {
						if !flg.NoLine {
							fmt.Println("--")
						}
					} else {
						if !flg.NoLine {
							fmt.Printf("%s:%s\n", l.Nr, l)
						} else {
							fmt.Println(l)
						}
					}
				}
			} else if flg.Hex {
				ctx.W = buffer.TermW

				for l := range buffer.Hex(ctx).Lines {
					fmt.Println(l)
				}
			} else {
				if h.Len() > 0 {
					fmt.Print(string(*h.MMap()))
				}
			}
		}
	})
}

func init() {
	flg := flags.Get()

	Fox.PersistentFlags().BoolVarP(&flg.Print, "print", "p", false, "print to console (no UI)")
	Fox.PersistentFlags().BoolVarP(&flg.NoFile, "no-file", "", false, "don't print filenames")
	Fox.PersistentFlags().BoolVarP(&flg.NoLine, "no-line", "", false, "don't print line numbers")

	Fox.Flags().BoolVarP(&flg.Hex, "hex", "x", false, "show file in canonical hex")

	Fox.Flags().BoolVarP(&flg.Limits.IsHead, "head", "h", false, "limit head of file by ...")
	Fox.Flags().BoolVarP(&flg.Limits.IsTail, "tail", "t", false, "limit tail of file by ...")
	Fox.Flags().IntVarP(&flg.Limits.Lines, "lines", "n", 0, "number of lines (default: 10)")
	Fox.Flags().IntVarP(&flg.Limits.Bytes, "bytes", "c", 0, "number of bytes (default: 16)")

	Fox.Flags().Lookup("lines").NoOptDefVal = "10"
	Fox.Flags().Lookup("bytes").NoOptDefVal = "16"

	Fox.Flags().VarP(&flg.Filters, "regexp", "e", "filter for lines that match pattern")
	Fox.Flags().IntVarP(&flg.Filters.Context, "context", "C", 0, "number of lines surrounding context of match")
	Fox.Flags().IntVarP(&flg.Filters.Before, "before", "B", 0, "number of lines leading context before match")
	Fox.Flags().IntVarP(&flg.Filters.After, "after", "A", 0, "number of lines trailing context after match")

	Fox.Flags().StringVarP(&flg.UI.State, "state", "", "", "sets the used UI state flags")
	Fox.Flags().StringVarP(&flg.UI.Theme, "theme", "", "", "sets the used UI theme")

	Fox.Flags().StringVarP(&flg.LLM.Model, "model", "", "", "sets the used Ollama model")

	Fox.Flags().StringVarP(&flg.Bag.Path, "file", "f", flags.BagName, "evidence bag file name")
	Fox.Flags().VarP(&flg.Bag.Mode, "mode", "m", "evidence bag file mode")
	Fox.Flags().StringVarP(&flg.Bag.Key, "key", "k", "", "key phrase for evidence bag signing with HMAC")
	Fox.Flags().StringVarP(&flg.Bag.Url, "url", "u", "", "url to also send evidence data too")
	Fox.Flags().BoolVarP(&flg.Bag.No, "no-bag", "", false, "don't write an evidence bag")

	Fox.Flags().Lookup("mode").NoOptDefVal = string(flags.BagModeText)

	Fox.Flags().BoolVarP(&flg.Opt.Raw, "raw", "r", false, "don't process files at all")
	Fox.Flags().BoolVarP(&flg.Opt.NoConvert, "no-convert", "", false, "don't convert automatically")
	Fox.Flags().BoolVarP(&flg.Opt.NoDeflate, "no-deflate", "", false, "don't deflate automatically")
	Fox.Flags().BoolVarP(&flg.Opt.NoPlugins, "no-plugins", "", false, "don't run autostart plugins")

	Fox.Flags().BoolVarP(&flg.Alias.Text, "text", "T", false, "short for --mode=text")
	Fox.Flags().BoolVarP(&flg.Alias.Json, "json", "j", false, "short for --mode=json")
	Fox.Flags().BoolVarP(&flg.Alias.Jsonl, "jsonl", "J", false, "short for --mode=jsonl")
	Fox.Flags().BoolVarP(&flg.Alias.Sqlite, "sqlite", "S", false, "short for --mode=sqlite")
	Fox.Flags().BoolVarP(&flg.Alias.Xml, "xml", "X", false, "short for --mode=xml")

	Fox.Flags().BoolP("help", "", false, "shows this message")
	Fox.Flags().BoolP("version", "", false, "shows the version")

	Fox.MarkFlagsMutuallyExclusive("head", "tail")

	Fox.SetHelpTemplate(fmt.Sprintf(fox.Usage, fox.Version))
	Fox.SetVersionTemplate(fmt.Sprintf("%s %s\n", fox.Product, fox.Version))

	Fox.AddCommand(sub.Counts)
	Fox.AddCommand(sub.Deflate)
	Fox.AddCommand(sub.Hash)
	Fox.AddCommand(sub.Strings)
}
