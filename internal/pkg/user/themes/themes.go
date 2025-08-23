package themes

import (
	"github.com/spf13/viper"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

const Filename = ".fox_themes"

type Themes struct {
	Themes map[string]Theme `mapstructure:"Theme"`
}

type Theme struct {
	Name     string
	Base     Style
	Surface0 Style
	Surface1 Style
	Surface2 Style
	Surface3 Style
	Overlay0 Style
	Overlay1 Style
	Subtext0 Style
	Subtext1 Style
	Subtext2 Style
}

type Style struct {
	Fg int32
	Bg int32
}

func New() *Themes {
	ts := new(Themes)

	cfg := viper.New()

	cfg.SetConfigName(Filename)
	cfg.SetConfigType("toml")
	cfg.AddConfigPath("$HOME")

	err := cfg.ReadInConfig()

	if err != nil {
		return nil
	}

	err = cfg.Unmarshal(ts)

	if err != nil {
		sys.Error(err)
		return nil
	}

	return ts
}
