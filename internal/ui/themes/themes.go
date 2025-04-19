package themes

import (
    "strings"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/gdamore/tcell/v2"
)

const (
    Default = "default"
)

var (
    // global styles
    Base     tcell.Style
    Surface0 tcell.Style
    Surface1 tcell.Style
    Surface2 tcell.Style
    Surface3 tcell.Style
    Overlay0 tcell.Style
    Overlay1 tcell.Style
    Subtext0 tcell.Style
    Subtext1 tcell.Style
    Colors   []tcell.Style
)

type Themes struct {
    palettes map[string]palette
    names []string
    index int
}

type palette []int32

func New(name string) *Themes {
    t := Themes{
        palettes: map[string]palette{
            "Default": palette {
                0x888888, 0x333333,
                0x333333, 0x333333,
                0xffffff, 0x333333,
                0xffffff, 0x333333,
                0xcccccc, 0x0f88cd,
                0x333333, 0xf8340c,
                0x333333, 0x88cd0f,
                0x888888, 0x333333,
                0x333333, 0x333333,
                0xffffff,
                0xffffff,
                0xffffff,
                0xffffff,
                0xffffff,
                0xffffff,
            },

            "Monokai": palette {
                0x7f8490, 0x222327,
                0x595f6f, 0x2c2e34,
                0xe2e2e3, 0x414550,
                0xe2e2e3, 0x2c2e34,
                0x222327, 0x85d3f2,
                0x222327, 0xff6077,
                0x222327, 0xa7df78,
                0x595f6f, 0x222327,
                0x2c2e34, 0x222327,
                0xfc5d7c,
                0xf39660,
                0xe7c664,
                0x9ed072,
                0x76cce0,
                0xb39df3,
            },

            "Catppuccin-Latte": palette {
                0x4c4f69, 0xeff1f5,
                0xacb0be, 0xccd0da,
                0x4c4f69, 0xbcc0cc,
                0x4c4f69, 0xccd0da,
                0xeff1f5, 0x1e66f5,
                0xeff1f5, 0xd20f39,
                0xeff1f5, 0x40a02b,
                0xacb0be, 0xeff1f5,
                0xccd0da, 0xeff1f5,
                0xd20f39,
                0xfe640b,
                0xdf8e1d,
                0x179299,
                0x04a5e5,
                0x209fb5,
            },

            "Catppuccin-Frappe": palette {
                0xa5adce, 0x303446,
                0x626880, 0x414559,
                0xc6d0f5, 0x51576d,
                0xc6d0f5, 0x414559,
                0x303446, 0x8caaee,
                0x303446, 0xe78284,
                0x303446, 0xa6d189,
                0x626880, 0x303446,
                0x414559, 0x303446,
                0xe78284,
                0xef9f76,
                0xe5c890,
                0x81c8be,
                0x99d1db,
                0x85c1dc,
            },

            "Catppuccin-Macchiato": palette {
                0xa5adcb, 0x24273a,
                0x5b6078, 0x363a4f,
                0xcad3f5, 0x494d64,
                0xcad3f5, 0x363a4f,
                0x24273a, 0x8aadf4,
                0x24273a, 0xed8796,
                0x24273a, 0xa6da95,
                0x5b6078, 0x24273a,
                0x363a4f, 0x24273a,
                0xed8796,
                0xf5a97f,
                0xeed49f,
                0x8bd5ca,
                0x91d7e3,
                0x7dc4e4,
            },

            "Catppuccin-Mocha": palette {
                0xa6adc8, 0x1e1e2e,
                0x585b70, 0x313244,
                0xcdd6f4, 0x45475a,
                0xcdd6f4, 0x313244,
                0x1e1e2e, 0x89b4fa,
                0x1e1e2e, 0xf38ba8,
                0x1e1e2e, 0xa6e3a1,
                0x585b70, 0x1e1e2e,
                0x313244, 0x1e1e2e,
                0xf38ba8,
                0xfab387,
                0xf9e2af,
                0x94e2d5,
                0x94e2d5,
                0x74c7ec,
            }, 

            "Ansi": palette {
                0xC0C0C0, 0x000000,
                0xffffff, 0x808080,
                0x000000, 0xC0C0C0,
                0xffffff, 0x808080,
                0xffffff, 0x000080,
                0x000000, 0x800000,
                0x000000, 0x008000,
                0x808080, 0x000000,
                0x808080, 0x000000,
                0xff0000,
                0xff00ff,
                0xffff00,
                0x00ffff,
                0x0000ff,
                0x00ff00,
            },

            "Matrix": palette {
                0x008f11, 0x0d0208,
                0x003b00, 0x0d0208,
                0x00ff41, 0x0d0208,
                0x00ff41, 0x0d0208,
                0x0d0208, 0x00ff41,
                0x0d0208, 0x00ff41,
                0x0d0208, 0x00ff41,
                0x003b00, 0x0d0208,
                0x0d0208, 0x0d0208,
                0x00ff41,
                0x00ff41,
                0x00ff41,
                0x00ff41,
                0x00ff41,
                0x00ff41,
            },

            "Monochrome": palette {
                0xffffff, 0x000000,
                0xffffff, 0x000000,
                0xffffff, 0x000000,
                0xffffff, 0x000000,
                0x000000, 0xffffff,
                0x000000, 0xffffff,
                0x000000, 0xffffff,
                0xffffff, 0x000000,
                0x000000, 0x000000,
                0xffffff,
                0xffffff,
                0xffffff,
                0xffffff,
                0xffffff,
                0xffffff,
            },
        },
        names: []string{
            "Default",
            "Monokai",
            "Catppuccin-Latte",
            "Catppuccin-Frappe",
            "Catppuccin-Macchiato",
            "Catppuccin-Mocha",
            "Ansi",
            "Matrix",
            "Monochrome",
        },
        index: 0,
    }

    t.Load(name)

    return &t
}

func (t *Themes) Cycle() string {
    t.index += 1
    t.index %= len(t.names)

    n := t.names[t.index]

    t.Load(n)

    return n
}

func (t *Themes) Load(name string) error {
    t.index = -1

    for i, n := range t.names {
        if strings.ToLower(n) == strings.ToLower(name) {
            t.index = i
            break
        }
    }

    if t.index == -1 {
        fx.Error("theme not found")

        t.index = 0
    }

    p := t.palettes[t.names[t.index]]

    Base = newStyle(p[0], p[1])
    Surface0 = newStyle(p[2], p[3])
    Surface1 = newStyle(p[4], p[5])
    Surface2 = newStyle(p[6], p[7])
    Surface3 = newStyle(p[8], p[9])
    Overlay0 = newStyle(p[10], p[11])
    Overlay1 = newStyle(p[12], p[13])
    Subtext0 = newStyle(p[14], p[15])
    Subtext1 = newStyle(p[16], p[17])

    Colors = Colors[:0] // reset

    for i := 18; i < 24; i++ {
        Colors = append(Colors, newStyle(p[i], p[1]))
    }

    return nil
}

func newStyle(fg, bg int32) tcell.Style {
    return tcell.StyleDefault.
        Foreground(tcell.NewHexColor(fg)).
        Background(tcell.NewHexColor(bg))
}
