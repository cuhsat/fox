package plugins

import (
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fx/pkg/fx/sys"
	"github.com/cuhsat/fx/pkg/fx/types"
	"github.com/cuhsat/fx/pkg/fx/types/heapset"
	"github.com/cuhsat/fx/pkg/fx/user"
)

const (
	filename = ".fx_plugins"
)

var (
	Value string
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

func (p *Plugins) Execute(hs *heapset.HeapSet, name string) (string, bool) {
	pl, ok := p.Plugins[name]

	if ok {
		_, h := hs.Heap()

		all := strings.Join(hs.Files(), " ")

		cmd := pl.Exec
		cmd = strings.ReplaceAll(cmd, "$?", Value)
		cmd = strings.ReplaceAll(cmd, "$+", h.Path)
		cmd = strings.ReplaceAll(cmd, "$*", all)

		hs.OpenFile(sys.Exec(cmd), pl.Name, types.Stdout)

		return pl.Name, true
	}

	return "", false
}
