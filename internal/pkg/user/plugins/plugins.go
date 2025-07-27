package plugins

import (
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/user"
)

const (
	Filename = ".fox_plugins"
)

var (
	Input chan string
)

type Func func(file sys.File, base string)

type Plugins struct {
	Autostart map[string]Plugin `toml:"Autostart"`
	Hotkey    map[string]Plugin `toml:"Hotkey"`
}

type Plugin struct {
	re *regexp.Regexp

	Name     string   `toml:"Name"`
	Prompt   string   `toml:"Prompt"`
	Pattern  string   `toml:"Pattern"`
	Commands []string `toml:"Commands"`
}

func New() *Plugins {
	Input = make(chan string)

	ps := new(Plugins)

	ok, path := user.File(Filename)

	if !ok {
		return nil
	}

	_, err := toml.DecodeFile(path, &ps)

	if err != nil {
		sys.Error(err)
		return nil
	}

	return ps
}

func Close() {
	close(Input)
}

func (ps *Plugins) Autostarts() []Plugin {
	r := make([]Plugin, len(ps.Autostart))

	for key := range ps.Autostart {
		p := ps.Autostart[key]
		p.re = regexp.MustCompile(p.Pattern)

		r = append(r, p)
	}

	return r
}

func (p *Plugin) Match(s string) bool {
	if p.re != nil {
		return p.re.MatchString(s)
	} else {
		return false
	}
}

func (p *Plugin) Execute(file, base string, fn Func) {
	var value string

	if len(p.Prompt) > 0 {
		value = <-Input // blocking call
	}

	r := strings.NewReplacer(
		"{{value}}", value,
		"{{file}}", sys.Persist(file),
		"{{base}}", sys.Persist(base),
	)

	cmds := make([]string, len(p.Commands))

	for _, cmd := range p.Commands {
		cmds = append(cmds, r.Replace(cmd))
	}

	fn(sys.Exec(cmds), base)
}
