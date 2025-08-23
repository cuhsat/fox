package themes

import (
	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user"
)

const (
	Filename = ".fox_themes"
)

type Themes struct {
	Themes map[string]Theme `toml:"Theme"`
}

type Theme struct {
	Name     string `toml:"name"`
	Base     Style  `toml:"base"`
	Surface0 Style  `toml:"surface0"`
	Surface1 Style  `toml:"surface1"`
	Surface2 Style  `toml:"surface2"`
	Surface3 Style  `toml:"surface3"`
	Overlay0 Style  `toml:"overlay0"`
	Overlay1 Style  `toml:"overlay1"`
	Subtext0 Style  `toml:"subtext0"`
	Subtext1 Style  `toml:"subtext1"`
	Subtext2 Style  `toml:"subtext2"`
}

type Style struct {
	Fg int32 `toml:"fg"`
	Bg int32 `toml:"bg"`
}

func New() *Themes {
	ts := new(Themes)

	is, p := user.File(Filename)

	if !is {
		return nil
	}

	_, err := toml.DecodeFile(p, &ts)

	if err != nil {
		sys.Error(err)
		return nil
	}

	return ts
}
