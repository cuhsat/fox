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
	filename = ".fox_plugins"
)

const (
	varBase   = "$BASE"
	varFile   = "$FILE"
	varFiles  = "$FILES"
	varInput  = "$INPUT"
	varParent = "$PARENT"
)

var (
	Input chan string
)

type Callback func(path, base, title string)

type Plugins struct {
	Autostarts map[string]Autostart `toml:"Autostart"`
	Plugins    map[string]Plugin    `toml:"Plugin"`
}

type Autostart struct {
	re *regexp.Regexp

	Name     string   `toml:"Name"`
	Pattern  string   `toml:"Pattern"`
	Commands []string `toml:"Commands"`
}

type Plugin struct {
	Name     string   `toml:"Name"`
	Prompt   string   `toml:"Prompt"`
	Commands []string `toml:"Commands"`
}

func (a *Autostart) Match(p string) bool {
	return a.re.MatchString(p)
}

func (a *Autostart) Execute(file, base string, hs []string) (string, string) {
	cmds := make([]string, len(a.Commands))

	for _, cmd := range a.Commands {
		cmds = append(cmds, expand(cmd, file, base, "", hs))

	}

	return sys.Exec(cmds), title(base, a.Name, "")
}

func (p *Plugin) Execute(file, base string, hs []string, fn Callback) {
	input := ""

	if len(p.Prompt) > 0 {
		input = expand(<-Input, file, base, "", hs)
	}

	cmds := make([]string, len(p.Commands))

	for _, cmd := range p.Commands {
		cmds = append(cmds, expand(cmd, file, base, input, hs))
	}

	fn(sys.Exec(cmds), base, title(base, p.Name, input))
}

func (ps *Plugins) Automatic() (as []Autostart) {
	keys := make([]string, 0)

	for k := range ps.Autostarts {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		as = append(as, ps.Autostarts[k])
	}

	return
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

	for k, v := range ps.Autostarts {
		v.re = regexp.MustCompile(v.Pattern)
		ps.Autostarts[k] = v
	}

	return ps
}

func Close() {
	close(Input)
}

func title(base, path, input string) string {
	var s string

	if len(input) > 0 {
		s = fmt.Sprintf(":%s", input)
	}

	return fmt.Sprintf("%s (%s%s)", base, path, s)
}

func expand(s, file, base, input string, hs []string) string {
	files := strings.Join(hs, " ")

	s = strings.ReplaceAll(s, varBase, base)
	s = strings.ReplaceAll(s, varFile, file)
	s = strings.ReplaceAll(s, varFiles, files)
	s = strings.ReplaceAll(s, varInput, input)
	s = strings.ReplaceAll(s, varParent, filepath.Dir(file))

	return s
}
