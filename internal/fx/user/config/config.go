package config

import (
    "bytes"
    "errors"
    "os"
    "path/filepath"

    "github.com/cuhsat/fx/internal/fx/sys"
    "github.com/BurntSushi/toml"
)

const (
    filename = ".fxrc"
)

type Config struct {
    Theme  string `toml:"Theme"`
    Follow bool   `toml:"Follow"`
    Line   bool   `toml:"Line"`
    Wrap   bool   `toml:"Wrap"`
}

func New() *Config {
    cfg := new(Config)

    dir, err := os.UserHomeDir()

    if err != nil {
        sys.Error(err)
        dir = "."
    }

    p := filepath.Join(dir, filename)

    _, err = os.Stat(p)

    if errors.Is(err, os.ErrNotExist) {
        return cfg
    } else if err != nil {
        sys.Error(err)
    }

    _, err = toml.DecodeFile(p, &cfg)

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

    dir, err := os.UserHomeDir()

    if err != nil {
        sys.Error(err)
        dir = "."
    }

    p := filepath.Join(dir, filename)

    err = os.WriteFile(p, buf.Bytes(), 0600)

    if err != nil {
        sys.Error(err)
    }
}
