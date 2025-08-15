package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/hiforensics/fox/cmd/sub"
	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/fox/ai"
	"github.com/hiforensics/fox/internal/fox/ai/agent"
	"github.com/hiforensics/fox/internal/fox/ui"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/text"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/buffer"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/types/heapset"
	"github.com/hiforensics/fox/internal/pkg/types/mode"
	"github.com/hiforensics/fox/internal/pkg/user/config"
)

// Fox usage
var Usage = fmt.Sprintf(fox.Fox+`
The Swiss Army Knife for examining text files (%s)

Usage:
  fox [COMMAND] [FLAG ...] [PATH ...]

Positional arguments:
  Path(s) to open or '-' for STDIN

Commands:
  counts                   display line and byte counts
  deflate                  deflate compressed files
  entropy                  display file entropy
  hash                     display file hash sums
  strings                  display ASCII and Unicode strings

Print:
  -p, --print              print directly to console
      --no-file            don't print filenames
      --no-line            don't print line numbers

Deflate:
      --pass=PASSWORD      decrypt with password (RAR, ZIP)

Hex display:
  -x, --hex                show file in canonical hex

File limits:
  -h, --head               limit head of file by ...
  -t, --tail               limit tail of file by ...
  -n, --lines[=NUMBER]     number of lines (default: 10)
  -c, --bytes[=NUMBER]     number of bytes (default: 16)

Line filter:
  -e, --regexp=PATTERN     filter for lines that match pattern
  -C, --context=NUMBER     number of lines surrounding context of match
  -B, --before=NUMBER      number of lines leading context before match
  -A, --after=NUMBER       number of lines trailing context after match

AI flags:
  -m, --model=MODEL        AI model for the agent to use
  -q, --query=QUERY        AI query for the agent to process

UI flags:
      --state={N|W|T|-}    sets the used UI state flags
      --theme=THEME        sets the used UI theme

Evidence bag:
  -f, --file=FILE          evidence bag file name (default: "evidence")
      --mode=MODE          evidence bag file mode (default: "raw")
                             NONE, RAW, TEST, JSON, JSONL, XML, SQLITE

  -s, --sign[=PHRASE]      sign evidence bag with (HMAC-)SHA256
  -u, --url=URL            url to also send evidence data too
      --no-bag             don't write an evidence bag

Disable:
  -R, --readonly           don't write any new files
  -r, --raw                don't process files at all
      --no-convert         don't convert automatically
      --no-deflate         don't deflate automatically
      --no-plugins         don't run autostart plugins

Aliases:
  -L, --logstash           short for --url=http://localhost:8080
  -T, --text               short for --mode=text
  -j, --json               short for --mode=json
  -J, --jsonl              short for --mode=jsonl
  -S, --sqlite             short for --mode=sqlite
  -X, --xml                short for --mode=xml

Standard:
      --help               shows this message
      --version            shows the version

Example: print matching lines
  $ fox -pe "John Doe" ./**/*.evtx

Example: print content hashes
  $ fox hash -pt sha1 files.zip

Example: print first sector in hex
  $ fox -pxhc=512 image.dd > mbr

Type "fox help COMMAND" for more help...
`, fox.Version)

// Displays files
var Fox = &cobra.Command{
	Use:     "fox",
	Short:   "The Swiss Army Knife for examining text files",
	Long:    "The Swiss Army Knife for examining text files",
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

		if flg.Opt.Readonly {
			flg.Opt.NoPlugins = true
			flg.Bag.No = true
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

		if flg.Alias.Logstash {
			flg.Bag.Url = flags.BagUrlLogstash
		}

		// explicit set UI mode
		if flg.Hex {
			flg.UI.Mode = mode.Hex
		}

		// implicit set UI mode
		if len(flg.Filters.Patterns) > 0 {
			flg.UI.Mode = mode.Grep
		}

		if len(flg.AI.Query) > 0 && !flg.Print {
			sys.Exit("query requires print")
		}

		if len(flg.UI.State) > 0 {
			re := regexp.MustCompile("[^-nwtNWT]+")

			flg.UI.State = re.ReplaceAllString(flg.UI.State, "")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !flags.Get().Print {
			ui.Start(args, types.None)
		} else if len(args) == 0 {
			fmt.Print(Usage)
			os.Exit(0)
		} else {
			print(args)
		}
	},
	SilenceUsage: true,
}

// Execute fox
func Execute() error {
	return Fox.Execute()
}

func print(args []string) {
	flg := flags.Get()

	if len(flg.AI.Query) > 0 {
		ai.Init(config.New().Model)

		if !ai.IsInit() {
			sys.Exit(ai.ErrNotAvailable.Error())
		}
	}

	hs := heapset.New(args)
	defer hs.ThrowAway()

	hs.Each(func(h *heap.Heap) {
		if h.Type != types.Stdin {
			ctx := buffer.NewContext(h)

			if !flg.NoFile {
				fmt.Println(text.Title(h.String(), buffer.TermW))
			}

			if len(flg.AI.Query) > 0 {
				agent.New().Ask(flg.AI.Query, h)
			} else if len(flg.Filters.Patterns) > 0 {
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

	Fox.Flags().BoolVarP(&flg.Print, "print", "p", false, "print directly to console")
	Fox.Flags().BoolVarP(&flg.NoFile, "no-file", "", false, "don't print filenames")
	Fox.Flags().BoolVarP(&flg.NoLine, "no-line", "", false, "don't print line numbers")

	Fox.PersistentFlags().StringVarP(&flg.Deflate.Pass, "pass", "", "", "decrypt with password")

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

	Fox.Flags().StringVarP(&flg.AI.Model, "model", "m", "", "AI model for the agent to use")
	Fox.Flags().StringVarP(&flg.AI.Query, "query", "q", "", "AI query for the agent to process")

	Fox.Flags().StringVarP(&flg.UI.State, "state", "", "", "sets the used UI state flags")
	Fox.Flags().StringVarP(&flg.UI.Theme, "theme", "", "", "sets the used UI theme")

	Fox.Flags().StringVarP(&flg.Bag.Path, "file", "f", flags.BagName, "evidence bag file name")
	Fox.Flags().VarP(&flg.Bag.Mode, "mode", "", "evidence bag file mode")
	Fox.Flags().StringVarP(&flg.Bag.Sign, "sign", "s", "", "sign evidence bag with (HMAC-)SHA256")
	Fox.Flags().StringVarP(&flg.Bag.Url, "url", "u", "", "url to also send evidence data too")
	Fox.Flags().BoolVarP(&flg.Bag.No, "no-bag", "", false, "don't write an evidence bag")

	Fox.Flags().Lookup("mode").NoOptDefVal = string(flags.BagModeText)
	Fox.Flags().Lookup("sign").NoOptDefVal = "-"

	Fox.Flags().BoolVarP(&flg.Opt.Raw, "raw", "r", false, "don't process files at all")
	Fox.Flags().BoolVarP(&flg.Opt.Readonly, "readonly", "R", false, "don't write any new files")
	Fox.Flags().BoolVarP(&flg.Opt.NoConvert, "no-convert", "", false, "don't convert automatically")
	Fox.Flags().BoolVarP(&flg.Opt.NoDeflate, "no-deflate", "", false, "don't deflate automatically")
	Fox.Flags().BoolVarP(&flg.Opt.NoPlugins, "no-plugins", "", false, "don't run autostart plugins")

	Fox.Flags().BoolVarP(&flg.Alias.Logstash, "logstash", "L", false, "short for --url=http://localhost:8080")
	Fox.Flags().BoolVarP(&flg.Alias.Text, "text", "T", false, "short for --mode=text")
	Fox.Flags().BoolVarP(&flg.Alias.Json, "json", "j", false, "short for --mode=json")
	Fox.Flags().BoolVarP(&flg.Alias.Jsonl, "jsonl", "J", false, "short for --mode=jsonl")
	Fox.Flags().BoolVarP(&flg.Alias.Sqlite, "sqlite", "S", false, "short for --mode=sqlite")
	Fox.Flags().BoolVarP(&flg.Alias.Xml, "xml", "X", false, "short for --mode=xml")

	Fox.PersistentFlags().BoolP("help", "", false, "shows this message")
	Fox.Flags().BoolP("version", "", false, "shows the version")

	Fox.MarkFlagsMutuallyExclusive("head", "tail")

	Fox.SetErrPrefix(sys.Prefix)
	Fox.SetHelpTemplate(Usage)
	Fox.SetVersionTemplate(fmt.Sprintf("%s %s\n", fox.Product, fox.Version))

	Fox.AddCommand(sub.Counts)
	Fox.AddCommand(sub.Deflate)
	Fox.AddCommand(sub.Entropy)
	Fox.AddCommand(sub.Hash)
	Fox.AddCommand(sub.Strings)
}
