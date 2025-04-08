package config

import (
    "errors"
    "os"
    "path/filepath"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/BurntSushi/toml"
)

const (
    File = ".curc"
)

type Config struct {
    UI struct {
        Theme string `toml:"Theme"`
        Line  bool   `toml:"Line"`
        Wrap  bool   `toml:"Wrap"`
    }
}

func Load() (c Config) {
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
        return // defaults
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

    return
}
