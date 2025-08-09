package flags

import (
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/mode"
)

type Flags struct {
	Print  bool
	NoFile bool
	NoLine bool

	Hex bool

	Limits  Limits
	Filters Filters

	// evidence bag
	Bag struct {
		Mode BagMode
		Path string
		Key  string
		Url  string
		No   bool
	}

	// optional flags
	Opt struct {
		Raw       bool
		NoConvert bool
		NoDeflate bool
		NoPlugins bool
	}

	// llm flags
	LLM struct {
		Model string
	}

	// ui flags
	UI struct {
		Theme  string
		State  string
		Mode   mode.Mode
		Invoke types.Invoke
	}

	// alias flags
	Alias struct {
		Text   bool
		Json   bool
		Jsonl  bool
		Sqlite bool
		Xml    bool
	}

	// deflate command
	Deflate struct {
		Path string
	}

	// hash command
	Hash struct {
		Algo HashAlgo
	}

	// strings command
	Strings struct {
		Min int
	}
}

var (
	flg *Flags = nil // singleton
)

func Get() *Flags {
	if flg == nil {
		flg = new(Flags)
		flg.UI.Mode = mode.Default
	}

	return flg
}
