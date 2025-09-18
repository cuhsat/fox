package flags

import (
	"regexp"

	"github.com/cuhsat/fox/internal/pkg/types/mode"
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
		Case string
		File string
		Mode BagMode
		Key  string
		Url  string
		Auth string
		Ecs  bool
		Hec  bool
	}

	// optional flags
	Opt struct {
		Raw       bool
		Readonly  bool
		NoConvert bool
		NoDeflate bool
		NoPlugins bool
		NoMouse   bool
	}

	// ai flags
	AI struct {
		Model  string
		Query  string
		NumCtx int
		Temp   float64
		TopP   float64
		TopK   int
		Seed   int
	}

	// ui flags
	UI struct {
		Theme  string
		State  string
		Space  int
		Legacy bool
		Mode   mode.Mode
	}

	// alias flags
	Alias struct {
		Logstash bool
		Splunk   bool
		Text     bool
		Json     bool
		Jsonl    bool
		Sqlite   bool
		Xml      bool
	}

	// compare command
	Compare struct {
		Git bool
	}

	// deflate command
	Deflate struct {
		List bool
		Path string
		Pass string
	}

	// entropy command
	Entropy struct {
		Min float64
		Max float64
	}

	// hash command
	Hash struct {
		Algo HashAlgo
	}

	// strings command
	Strings struct {
		Ioc   bool
		Re    *regexp.Regexp
		Min   int
		Max   int
		Ascii bool
	}

	Credits bool
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
