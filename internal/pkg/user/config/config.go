package config

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user"
)

const (
	Filename = ".foxrc"
)

type Config struct {
	Model   string `toml:"Model"`
	Theme   string `toml:"Theme"`
	Follow  bool   `toml:"Follow"`
	Numbers bool   `toml:"Numbers"`
	Wrap    bool   `toml:"Wrap"`
}

func New() *Config {
	cfg := new(Config)

	is, p := user.File(Filename)

	if !is {
		return cfg
	}

	_, err := toml.DecodeFile(p, &cfg)

	if err != nil {
		sys.Error(err)
	}

	// higher ranking environment variables
	m := os.Getenv("FOX_MODEL")

	if len(m) > 0 {
		cfg.Model = m
	}

	t := os.Getenv("FOX_THEME")

	if len(t) > 0 {
		cfg.Theme = t
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

	_, p := user.File(Filename)

	err = os.WriteFile(p, buf.Bytes(), 0600)

	if err != nil {
		sys.Error(err)
	}
}
