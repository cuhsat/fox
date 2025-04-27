package themes

import (
    "errors"
    "os"
    "path/filepath"

    "github.com/cuhsat/fx/internal/fx/sys"
    "github.com/BurntSushi/toml"
)

const (
    filename = ".fx_themes"
)

type Themes struct {
    Themes map[string]Theme `toml:"Theme"`
}

type Theme struct {
    Name     string `toml:"Name"`
    Base     Style  `toml:"Base"`
    Surface0 Style  `toml:"Surface0"`
    Surface1 Style  `toml:"Surface1"`
    Surface2 Style  `toml:"Surface2"`
    Surface3 Style  `toml:"Surface3"`
    Overlay0 Style  `toml:"Overlay0"`
    Overlay1 Style  `toml:"Overlay1"`
    Subtext0 Style  `toml:"Subtext0"`
    Subtext1 Style  `toml:"Subtext1"`
    Subtext2 Style  `toml:"Subtext2"`
}

type Style struct {
    Fg int32 `toml:"Fg"`
    Bg int32 `toml:"Bg"`
}

func New() *Themes {
    ts := new(Themes)

    dir, err := os.UserHomeDir()

    if err != nil {
        sys.Error(err)
        dir = "."
    }

    p := filepath.Join(dir, filename)

    _, err = os.Stat(p)

    if errors.Is(err, os.ErrNotExist) {
        return nil
    } else if err != nil {
        sys.Error(err)
        return nil
    }

    _, err = toml.DecodeFile(p, &ts)

    if err != nil {
        sys.Error(err)
        return nil
    }

    return ts
}
