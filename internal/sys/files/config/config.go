package config

import (
    "bytes"
    "errors"
    "os"
    "path/filepath"

    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/cuhsat/fx/internal/sys"
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

func Load() *Config {
    var c Config

    // defaults
    c.Theme = themes.Default
    c.Follow = false
    c.Line = true
    c.Wrap = false

    dir, err := os.UserHomeDir()

    if err != nil {
        sys.Fatal(err)
    }

    p := filepath.Join(dir, filename)

    _, err = os.Stat(p)

    if errors.Is(err, os.ErrNotExist) {
        return &c // defaults
    } else if err != nil {
        sys.Fatal(err)
    }

    _, err = toml.DecodeFile(p, &c)

    if err != nil {
        sys.Fatal(err)
    }

    // higher ranking variables
    env := os.Getenv("FX_THEME")

    if len(env) > 0 {
        c.Theme = env
    }

    return &c
}

func (c *Config) Save() {
    b := new(bytes.Buffer)

    e := toml.NewEncoder(b)
    
    e.Indent = "" // no indent

    err := e.Encode(c)

    if err != nil {
        sys.Fatal(err)
    }

    dir, err := os.UserHomeDir()

    if err != nil {
        sys.Fatal(err)
    }

    p := filepath.Join(dir, filename)

    err = os.WriteFile(p, b.Bytes(), 0600)

    if err != nil {
        sys.Fatal(err)
    }
}
