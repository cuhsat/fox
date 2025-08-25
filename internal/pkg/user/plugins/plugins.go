package plugins

import (
	"regexp"
	"strings"

	"github.com/spf13/viper"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user"
)

var Input chan string

type Callback func(path, base, dir string)

type Plugins struct {
	Auto   map[string]Plugin `mapstructure:"auto"`
	Hotkey map[string]Plugin `mapstructure:"hotkey"`
}

type Plugin struct {
	re *regexp.Regexp

	Name string
	Mode string
	Mask string
	Exec []string
}

func New() *Plugins {
	Input = make(chan string)

	ps := new(Plugins)

	cfg := viper.New()

	if !user.LoadConfig(cfg, "plugins") {
		return nil
	}

	err := cfg.Unmarshal(ps)

	if err != nil {
		sys.Error(err)
		return nil
	}

	return ps
}

func Close() {
	close(Input)
}

func (ps *Plugins) Autos() []Plugin {
	as := make([]Plugin, len(ps.Auto))

	for key := range ps.Auto {
		p := ps.Auto[key]
		p.re = regexp.MustCompile(p.Mask)

		as = append(as, p)
	}

	return as
}

func (p *Plugin) Match(mask string) bool {
	if p.re != nil {
		return p.re.MatchString(mask)
	} else {
		return false
	}
}

func (p *Plugin) Execute(file, base string, fn Callback) {
	var val, temp string

	// blocking call
	if len(p.Mode) > 0 {
		val = <-Input
	}

	// create temp dir if necessary
	for _, cmd := range p.Exec {
		if strings.Contains(cmd, "$TEMP") {
			temp = user.TempDir("plugin")
			break
		}
	}

	// replace and persist
	rep := strings.NewReplacer(
		"$BASE", user.Persist(base),
		"$FILE", user.Persist(file),
		"$TEMP", temp,
		"$INPUT", val,
	)

	cmds := make([]string, len(p.Exec))

	for _, cmd := range p.Exec {
		cmds = append(cmds, rep.Replace(cmd))
	}

	fn(sys.Call(cmds).Name(), base, temp)
}
