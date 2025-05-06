package plugins

import (
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/heapset"
	"github.com/cuhsat/fx/internal/pkg/user"
)

const (
	filename = ".fx_plugins"
)

var (
	Input chan string
)

type Plugins struct {
	Autostart map[string][]Autostart `toml:"Autostart"`
	Plugins   map[string]Plugin      `toml:"Plugin"`
}

type Autostart struct {
	Path string `toml:"Path"`
	Exec string `toml:"Exec"`
}

type Plugin struct {
	Name  string `toml:"Name"`
	Exec  string `toml:"Exec"`
	Input string `toml:"Input"`
}

func New() *Plugins {
	Input = make(chan string)

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

func Close() {
	close(Input)
}

func (p *Plugin) Autostart(path string) {

}

func (p *Plugin) Execute(hs *heapset.HeapSet, fn func()) {
	var v string

	if len(p.Input) > 0 {
		v = <-Input
	}

	_, h := hs.Heap()

	all := strings.Join(hs.Files(), " ")

	cmd := p.Exec
	cmd = strings.ReplaceAll(cmd, "$?", v)
	cmd = strings.ReplaceAll(cmd, "$+", h.Path)
	cmd = strings.ReplaceAll(cmd, "$*", all)

	hs.OpenFile(sys.Exec(cmd), p.Name, types.Stdout)

	fn()
}
