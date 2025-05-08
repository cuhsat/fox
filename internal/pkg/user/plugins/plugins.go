package plugins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/user"
)

const (
	filename = ".fx_plugins"
)

const (
	varBase  = "$BASE"
	varFile  = "$FILE"
	varFiles = "$FILES"
	varInput = "$INPUT"
)

const (
	ext = "text"
)

var (
	Input chan string
)

type Callback func(p, t string)

type Plugins struct {
	Starts  map[string][]Autostart `toml:"Plugins"`
	Plugins map[string]Plugin      `toml:"Plugin"`
}

type Autostart struct {
	re *regexp.Regexp

	Name    string `toml:"Name"`
	Pattern string `toml:"Pattern"`
	Command string `toml:"Command"`
	Output  string `toml:"Output"`
}

type Plugin struct {
	Name    string `toml:"Name"`
	Prompt  string `toml:"Prompt"`
	Command string `toml:"Command"`
	Output  string `toml:"Output"`
}

func (a *Autostart) Match(p string) bool {
	return a.re.MatchString(p)
}

func (a *Autostart) Execute(f, b string) string {
	var e string = ext

	if len(a.Output) > 0 {
		e = a.Output
	}

	cmd := a.Command
	cmd = strings.ReplaceAll(cmd, varBase, b)
	cmd = strings.ReplaceAll(cmd, varFile, f)

	return sys.Exec(cmd, e)
}

func (p *Plugin) Execute(f, b string, hs []string, fn Callback) {
	var e string = ext

	if len(p.Output) > 0 {
		e = p.Output
	}

	var s string

	if len(p.Prompt) > 0 {
		s = <-Input
	}

	fs := strings.Join(hs, " ")

	cmd := p.Command
	cmd = strings.ReplaceAll(cmd, varBase, b)
	cmd = strings.ReplaceAll(cmd, varFile, f)
	cmd = strings.ReplaceAll(cmd, varFiles, fs)
	cmd = strings.ReplaceAll(cmd, varInput, s)

	fn(sys.Exec(cmd, e), fmt.Sprintf("%s/%s", b, s))
}

func (ps *Plugins) Autostarts() []Autostart {
	return ps.Starts["Autostart"]
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

	for _, ps := range ps.Starts {
		for i, p := range ps {
			ps[i].re = regexp.MustCompile(p.Pattern)
		}
	}

	return ps
}

func Close() {
	close(Input)
}
