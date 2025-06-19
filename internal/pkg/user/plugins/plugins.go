package plugins

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user"
)

const (
	Filename = ".fox_plugins"
)

var (
	Input chan string
)

type Callback func(path, base, title string)

type Plugins struct {
	Automatics map[string]Plugin `toml:"Autostart"`
	Shortcuts  map[string]Plugin `toml:"Plugin"`
}

type Plugin struct {
	re *regexp.Regexp

	Name     string   `toml:"Name"`
	Prompt   string   `toml:"Prompt"`
	Pattern  string   `toml:"Pattern"`
	Commands []string `toml:"Commands"`
}

func (p *Plugin) Match(s string) bool {
	return p.re.MatchString(s)
}

func (p *Plugin) Execute(file, base string, hs []string, fn Callback) (string, string) {
	var input string
	var title string

	if len(p.Prompt) > 0 {
		input = <-Input
	}

	title = fmt.Sprintf("%s (%s%s)", base, p.Name, title)

	cs := make([]string, len(p.Commands))

	for _, cmd := range p.Commands {
		cmd = strings.ReplaceAll(cmd, "$BASE", base)
		cmd = strings.ReplaceAll(cmd, "$FILE", file)
		cmd = strings.ReplaceAll(cmd, "$FILES", strings.Join(hs, " "))
		cmd = strings.ReplaceAll(cmd, "$INPUT", input)
		cmd = strings.ReplaceAll(cmd, "$FOLDER", filepath.Dir(file))

		cs = append(cs, cmd)
	}

	if fn != nil {
		fn(sys.Exec(cs), base, title)
	} else {
		return sys.Exec(cs), title
	}

	return "", ""
}

func (ps *Plugins) Auto() (as []Plugin) {
	keys := make([]string, 0)

	for k := range ps.Automatics {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		p := ps.Automatics[k]
		p.re = regexp.MustCompile(p.Pattern)
		as = append(as, p)
	}

	return
}

func New() *Plugins {
	Input = make(chan string)

	ps := new(Plugins)

	ok, p := user.File(Filename)

	if !ok {
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
