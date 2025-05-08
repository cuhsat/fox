package config

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/user"
)

const (
	filename = ".fxrc"
)

type Config struct {
	Theme string `toml:"Theme"`
	Tail  bool   `toml:"Tail"`
	Line  bool   `toml:"Line"`
	Wrap  bool   `toml:"Wrap"`
}

func New() *Config {
	cfg := new(Config)

	is, p := user.Config(filename)

	if !is {
		return cfg
	}

	_, err := toml.DecodeFile(p, &cfg)

	if err != nil {
		sys.Error(err)
	}

	// higher ranking variables
	env := os.Getenv("FX_THEME")

	if len(env) > 0 {
		cfg.Theme = env
	}

	return cfg
}

func (cfg *Config) Save() {
	buf := new(bytes.Buffer)

	enc := toml.NewEncoder(buf)
	enc.Indent = "" // no indent

	err := enc.Encode(cfg)

	if err != nil {
		sys.Error(err)
		return
	}

	_, p := user.Config(filename)

	err = os.WriteFile(p, buf.Bytes(), 0600)

	if err != nil {
		sys.Error(err)
	}
}
