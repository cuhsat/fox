package config

import (
    "errors"
    "os"
    "path/filepath"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/BurntSushi/toml"
)

const (
    File = ".curc"
)

type Config struct {
    CU struct {
        Hash string `toml:"Hash"`
    }
    UI struct {
        Theme string `toml:"Theme"`
        Line  bool   `toml:"Line"`
        Wrap  bool   `toml:"Wrap"`
    }
}

// singleton
var instance *Config = nil

func NewConfig() *Config {
    if instance == nil {
        instance = load();
    }

    return instance;
}

func load() *Config {
    var c Config

    // defaults CU
    c.CU.Hash = heap.Sha256

    // defaults UI
    c.UI.Theme = theme.Default
    c.UI.Line = true
    c.UI.Wrap = false

    dir, err := os.UserHomeDir()

    if err != nil {
        fs.Panic(err)
    }

    f := filepath.Join(dir, File)

    _, err = os.Stat(f)

    if errors.Is(err, os.ErrNotExist) {
        return &c // defaults
    } else if err != nil {
        fs.Panic(err)
    }

    _, err = toml.DecodeFile(f, &c)

    if err != nil {
        fs.Panic(err)
    }

    env := os.Getenv("CU_THEME")

    if len(env) > 0 {
        c.UI.Theme = env
    }

    return &c
}
