package themes

import (
	"github.com/spf13/viper"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user"
)

type Themes struct {
	Themes map[string]struct {
		Name     string
		Terminal [2]int32
		Surface0 [2]int32
		Surface1 [2]int32
		Surface2 [2]int32
		Surface3 [2]int32
		Overlay0 [2]int32
		Overlay1 [2]int32
		Subtext0 [2]int32
		Subtext1 [2]int32
		Subtext2 [2]int32
	} `mapstructure:"Theme"`
}

func New() *Themes {
	ts := new(Themes)

	cfg := viper.New()

	if !user.LoadConfig(cfg, "themes") {
		return nil
	}

	err := cfg.Unmarshal(ts)

	if err != nil {
		sys.Error(err)
		return nil
	}

	return ts
}
