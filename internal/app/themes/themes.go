package themes

import (
    "strings"

    "github.com/cuhsat/cu/internal/sys"
    "github.com/gdamore/tcell/v2"
)

const (
    Default = "monokai"
)

// global styles
var (
    Output tcell.Style
    Header tcell.Style
    Input  tcell.Style
    Error  tcell.Style
    Mode   tcell.Style
    Hint   tcell.Style
    Rule   tcell.Style
    Info   tcell.Style
    Line   tcell.Style
    Colors []tcell.Style
)

type Themes struct {
    palettes map[string][]int32
    names []string
    index int
}

func NewThemes(name string) *Themes {
    t := Themes{
        palettes: map[string][]int32{
            "Monokai": []int32 {
                0x7f8490, // output foreground
                0x222327, // output background
                0xe2e2e3, // header foreground
                0x2c2e34, // header background
                0xe2e2e3, // input foreground
                0x414550, // input background
                0x222327, // error foreground
                0xff6077, // error background
                0x222327, // mode foreground
                0xa7df78, // mode background
                0x595f6f, // hint foreground
                0x222327, // hint background
                0x2c2e34, // rule foreground
                0x222327, // rule background
                0x222327, // info foreground
                0x85d3f2, // info background
                0x595f6f, // line foreground
                0x2c2e34, // line background
                0xfc5d7c, // highlight 1
                0xf39660, // highlight 2
                0xe7c664, // highlight 3
                0x9ed072, // highlight 4
                0x76cce0, // highlight 5
                0xb39df3, // highlight 6
            },

            "Catppuccin-Latte": []int32 {
                0x4c4f69, // Subtext0
                0xeff1f5, // Base
                0x4c4f69, // Text
                0xccd0da, // Surface0
                0x4c4f69, // Text
                0xbcc0cc, // Surface1
                0xeff1f5, // Base
                0xd20f39, // Red
                0xeff1f5, // Base
                0x40a02b, // Green
                0xacb0be, // Surface2
                0xeff1f5, // Base
                0xccd0da, // Surface0
                0xeff1f5, // Base
                0xeff1f5, // Base
                0x1e66f5, // Blue
                0xacb0be, // Surface2
                0xccd0da, // Surface0
                0xd20f39, // Red
                0xfe640b, // Peach
                0xdf8e1d, // Yellow
                0x179299, // Teal
                0x04a5e5, // Sky
                0x209fb5, // Sapphire
            },

            "Catppuccin-Frappe": []int32 {
                0xa5adce, // Subtext0
                0x303446, // Base
                0xc6d0f5, // Text
                0x414559, // Surface0
                0xc6d0f5, // Text
                0x51576d, // Surface1
                0x303446, // Base
                0xe78284, // Red
                0x303446, // Base
                0xa6d189, // Green
                0x626880, // Surface2
                0x303446, // Base
                0x414559, // Surface0
                0x303446, // Base
                0x303446, // Base
                0x8caaee, // Blue
                0x626880, // Surface2
                0x414559, // Surface0
                0xe78284, // Red
                0xef9f76, // Peach
                0xe5c890, // Yellow
                0x81c8be, // Teal
                0x99d1db, // Sky
                0x85c1dc, // Sapphire
            },

            "Catppuccin-Macchiato": []int32 {
                0xa5adcb, // Subtext0
                0x24273a, // Base
                0xcad3f5, // Text
                0x363a4f, // Surface0
                0xcad3f5, // Text
                0x494d64, // Surface1
                0x24273a, // Base
                0xed8796, // Red
                0x24273a, // Base
                0xa6da95, // Green
                0x5b6078, // Surface2
                0x24273a, // Base
                0x363a4f, // Surface0
                0x24273a, // Base
                0x24273a, // Base
                0x8aadf4, // Blue
                0x5b6078, // Surface2
                0x363a4f, // Surface0
                0xed8796, // Red
                0xf5a97f, // Peach
                0xeed49f, // Yellow
                0x8bd5ca, // Teal
                0x91d7e3, // Sky
                0x7dc4e4, // Sapphire
            },

            "Catppuccin-Mocha": []int32 {
                0xa6adc8, // Subtext0
                0x1e1e2e, // Base
                0xcdd6f4, // Text
                0x313244, // Surface0
                0xcdd6f4, // Text
                0x45475a, // Surface1
                0x1e1e2e, // Base
                0xf38ba8, // Red
                0x1e1e2e, // Base
                0xa6e3a1, // Green
                0x585b70, // Surface2
                0x1e1e2e, // Base
                0x313244, // Surface0
                0x1e1e2e, // Base
                0x1e1e2e, // Base
                0x89b4fa, // Blue
                0x585b70, // Surface2
                0x313244, // Surface0
                0xf38ba8, // Red
                0xfab387, // Peach
                0xf9e2af, // Yellow
                0x94e2d5, // Teal
                0x94e2d5, // Sky
                0x74c7ec, // Sapphire
            },    

            "Matrix": []int32 {
                0x008f11, // output foreground
                0x0d0208, // output background
                0x00ff41, // header foreground
                0x0d0208, // header background
                0x00ff41, // input foreground
                0x0d0208, // input background
                0x0d0208, // error foreground
                0x00ff41, // error background
                0x0d0208, // mode foreground
                0x00ff41, // mode background
                0x003b00, // hint foreground
                0x0d0208, // hint background
                0x0d0208, // rule foreground
                0x0d0208, // rule background
                0x0d0208, // info foreground
                0x00ff41, // info background
                0x003b00, // line foreground
                0x0d0208, // line background
                0x00ff41, // highlight 1
                0x00ff41, // highlight 2
                0x00ff41, // highlight 3
                0x00ff41, // highlight 4
                0x00ff41, // highlight 5
                0x00ff41, // highlight 6
            },

            "Monochrome": []int32 {
                0xffffff, // output foreground
                0x000000, // output background
                0xffffff, // header foreground
                0x000000, // header background
                0xffffff, // input foreground
                0x000000, // input background
                0x000000, // error foreground
                0xffffff, // error background
                0x000000, // mode foreground
                0xffffff, // mode background
                0xffffff, // hint foreground
                0x000000, // hint background
                0x000000, // rule foreground
                0x000000, // rule background
                0x000000, // info foreground
                0xffffff, // info background
                0xffffff, // line foreground
                0x000000, // line background
                0xffffff, // highlight 1
                0xffffff, // highlight 2
                0xffffff, // highlight 3
                0xffffff, // highlight 4
                0xffffff, // highlight 5
                0xffffff, // highlight 6
            },
        },
        names: []string{
            "Monokai",
            "Catppuccin-Latte",
            "Catppuccin-Frappe",
            "Catppuccin-Macchiato",
            "Catppuccin-Mocha",
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

func (t *Themes) Load(name string) {
    t.index = -1

    for i, n := range t.names {
        if strings.ToLower(n) == strings.ToLower(name) {
            t.index = i
            break
        }
    }

    if t.index == -1 {
        sys.Fatal("theme not found")
    }

    p := t.palettes[t.names[t.index]]

    Output = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[0])).
        Background(tcell.NewHexColor(p[1]))

    Header = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[2])).
        Background(tcell.NewHexColor(p[3]))

    Input = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[4])).
        Background(tcell.NewHexColor(p[5]))

    Error = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[6])).
        Background(tcell.NewHexColor(p[7]))

    Mode = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[8])).
        Background(tcell.NewHexColor(p[9]))

    Hint = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[10])).
        Background(tcell.NewHexColor(p[11]))

    Rule = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[12])).
        Background(tcell.NewHexColor(p[13]))

    Info = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[14])).
        Background(tcell.NewHexColor(p[15]))

    Line = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[16])).
        Background(tcell.NewHexColor(p[17]))

    Colors = Colors[:0] // reset

    for i := 18; i < 24; i++ {
        Colors = append(Colors, tcell.StyleDefault.
            Foreground(tcell.NewHexColor(p[i])).
            Background(tcell.NewHexColor(p[1])))
    }
}
