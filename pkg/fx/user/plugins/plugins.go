package plugins

import (
	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fx/pkg/fx/sys"
	"github.com/cuhsat/fx/pkg/fx/user"
)

const (
	filename = ".fx_plugins"
)

type Plugins struct {
	Plugins map[string]Plugin `toml:"Plugin"`
}

type Plugin struct {
	Name string `toml:"Name"`
	Exec string `toml:"Exec"`
}

func New() *Plugins {
	ps := new(Plugins)

	is, p := user.Config(filename)

	if !is {
		return nil
	}

	_, err := toml.DecodeFile(p, &ps)

	if err != nil {
		sys.Error(err)
		return nil
	}

	return ps
}
