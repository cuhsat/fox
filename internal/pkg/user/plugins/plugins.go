package plugins

import (
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/user"
)

const (
	filename = ".fx_plugins"
)

var (
	Input chan string
)

type Callback func(p, t string)

type Plugins struct {
	Starts  map[string][]Start `toml:"Plugins"`
	Plugins map[string]Plugin  `toml:"Plugin"`
}

type Start struct {
	Name string `toml:"Name"`
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

func (s *Start) Execute(f, b string) string {
	cmd := s.Exec
	cmd = strings.ReplaceAll(cmd, "$.", b)
	cmd = strings.ReplaceAll(cmd, "$+", f)

	return sys.Exec(cmd)
}

func (p *Plugin) Execute(f, b string, hs []string, fn Callback) {
	var v string

	if len(p.Input) > 0 {
		v = <-Input
	}

	cmd := p.Exec
	cmd = strings.ReplaceAll(cmd, "$?", v)
	cmd = strings.ReplaceAll(cmd, "$.", b)
	cmd = strings.ReplaceAll(cmd, "$+", f)
	cmd = strings.ReplaceAll(cmd, "$*", strings.Join(hs, " "))

	fn(sys.Exec(cmd), p.Name)
}
