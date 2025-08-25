package themes

import (
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user/themes"
)

const (
	Default = "examiner-dark"
)

var (
	Cursor   tcell.Color
	Terminal tcell.Style
	Surface0 tcell.Style
	Surface1 tcell.Style
	Surface2 tcell.Style
	Surface3 tcell.Style
	Overlay0 tcell.Style
	Overlay1 tcell.Style
	Subtext0 tcell.Style
	Subtext1 tcell.Style
	Subtext2 tcell.Style
)

type Themes struct {
	palettes map[string]palette
	names    []string
	index    int
}

type palette []int32

func New(name string) *Themes {
	t := Themes{
		palettes: map[string]palette{
			"Examiner-Dark": {
				0x777777, 0x111111, // Terminal
				0x333333, 0x333333, // Surface0
				0xeeeeee, 0x333333, // Surface1
				0xeeeeee, 0x333333, // Surface2
				0xeeeeee, 0x0f88cd, // Surface3
				0xeeeeee, 0xf8340c, // Overlay0
				0x333333, 0x88cd0f, // Overlay1
				0x333333, 0x111111, // Subtext0
				0x333333, 0x111111, // Subtext1
				0xeeeeee, 0x111111, // Subtext2
			},

			"Examiner-Light": {
				0x888888, 0xeeeeee, // Terminal
				0xcccccc, 0xcccccc, // Surface0
				0x111111, 0xcccccc, // Surface1
				0x111111, 0xcccccc, // Surface2
				0xeeeeee, 0x0f88cd, // Surface3
				0xeeeeee, 0xf8340c, // Overlay0
				0x333333, 0x88cd0f, // Overlay1
				0xcccccc, 0xeeeeee, // Subtext0
				0xcccccc, 0xeeeeee, // Subtext1
				0x111111, 0xeeeeee, // Subtext2
			},

			"Catppuccin-Latte": {
				0x4c4f69, 0xeff1f5, // Terminal
				0xacb0be, 0xccd0da, // Surface0
				0x4c4f69, 0xbcc0cc, // Surface1
				0x4c4f69, 0xccd0da, // Surface2
				0xeff1f5, 0x1e66f5, // Surface3
				0xeff1f5, 0xd20f39, // Overlay0
				0xeff1f5, 0x40a02b, // Overlay1
				0xacb0be, 0xeff1f5, // Subtext0
				0xccd0da, 0xeff1f5, // Subtext1
				0xd20f39, 0xeff1f5, // Subtext2
			},

			"Catppuccin-Frappe": {
				0xa5adce, 0x303446, // Terminal
				0x626880, 0x414559, // Surface0
				0xc6d0f5, 0x51576d, // Surface1
				0xc6d0f5, 0x414559, // Surface2
				0x303446, 0x8caaee, // Surface3
				0x303446, 0xe78284, // Overlay0
				0x303446, 0xa6d189, // Overlay1
				0x626880, 0x303446, // Subtext0
				0x414559, 0x303446, // Subtext1
				0xe78284, 0x303446, // Subtext2
			},

			"Catppuccin-Macchiato": {
				0xa5adcb, 0x24273a, // Terminal
				0x5b6078, 0x363a4f, // Surface0
				0xcad3f5, 0x494d64, // Surface1
				0xcad3f5, 0x363a4f, // Surface2
				0x24273a, 0x8aadf4, // Surface3
				0x24273a, 0xed8796, // Overlay0
				0x24273a, 0xa6da95, // Overlay1
				0x5b6078, 0x24273a, // Subtext0
				0x363a4f, 0x24273a, // Subtext1
				0xed8796, 0x24273a, // Subtext2
			},

			"Catppuccin-Mocha": {
				0xa6adc8, 0x1e1e2e, // Terminal
				0x585b70, 0x313244, // Surface0
				0xcdd6f4, 0x45475a, // Surface1
				0xcdd6f4, 0x313244, // Surface2
				0x1e1e2e, 0x89b4fa, // Surface3
				0x1e1e2e, 0xf38ba8, // Overlay0
				0x1e1e2e, 0xa6e3a1, // Overlay1
				0x585b70, 0x1e1e2e, // Subtext0
				0x313244, 0x1e1e2e, // Subtext1
				0xf38ba8, 0x1e1e2e, // Subtext2
			},

			"Solarized-Dark": {
				0x93a1a1, 0x002b36, // Terminal
				0x073642, 0x073642, // Surface0
				0xfdf6e3, 0x073642, // Surface1
				0xfdf6e3, 0x073642, // Surface2
				0xfdf6e3, 0x586e75, // Surface3
				0xfdf6e3, 0xdc322f, // Overlay0
				0xfdf6e3, 0x859900, // Overlay1
				0x586e75, 0x002b36, // Subtext0
				0x073642, 0x002b36, // Subtext1
				0xcb4b16, 0x002b36, // Subtext2
			},

			"Solarized-Light": {
				0x586e75, 0xfdf6e3, // Terminal
				0xeee8d5, 0xeee8d5, // Surface0
				0x002b36, 0xeee8d5, // Surface1
				0x002b36, 0xeee8d5, // Surface2
				0xfdf6e3, 0x93a1a1, // Surface3
				0xfdf6e3, 0xdc322f, // Overlay0
				0xfdf6e3, 0x859900, // Overlay1
				0x93a1a1, 0xfdf6e3, // Subtext0
				0xeee8d5, 0xfdf6e3, // Subtext1
				0xb58900, 0xfdf6e3, // Subtext2
			},

			"VSCode-Dark": {
				0xdee1e6, 0x282828, // Terminal
				0xdee1e6, 0x313131, // Surface0
				0xdee1e6, 0x444444, // Surface1
				0xdee1e6, 0x313131, // Surface2
				0x1a1a1a, 0x569cd6, // Surface3
				0x1a1a1a, 0xd16969, // Overlay0
				0x1a1a1a, 0xb5cea8, // Overlay1
				0x626262, 0x282828, // Subtext0
				0x313131, 0x282828, // Subtext1
				0xd3967d, 0x282828, // Subtext2
			},

			"VSCode-Light": {
				0x343434, 0xe7e7e7, // Terminal
				0x343434, 0xdfdfdf, // Surface0
				0x343434, 0xcfcfcf, // Surface1
				0x343434, 0xdfdfdf, // Surface2
				0xe7e7e7, 0x007acc, // Surface3
				0xe7e7e7, 0xff0000, // Overlay0
				0xe7e7e7, 0x008000, // Overlay1
				0xafafaf, 0xe7e7e7, // Subtext0
				0xafafaf, 0xe7e7e7, // Subtext1
				0xc72e0f, 0xe7e7e7, // Subtext2
			},

			"Monokai": {
				0x7f8490, 0x222327, // Terminal
				0x595f6f, 0x2c2e34, // Surface0
				0xe2e2e3, 0x414550, // Surface1
				0xe2e2e3, 0x2c2e34, // Surface2
				0x222327, 0xa7df78, // Surface3
				0x222327, 0xff6077, // Overlay0
				0x222327, 0x85d3f2, // Overlay1
				0x595f6f, 0x222327, // Subtext0
				0x2c2e34, 0x222327, // Subtext1
				0xf39660, 0x222327, // Subtext2
			},

			"Darcula": {
				0x727272, 0x2b2b2b, // Terminal
				0x393939, 0x393939, // Surface0
				0x727272, 0x393939, // Surface1
				0x727272, 0x393939, // Surface2
				0x2b2b2b, 0x727272, // Surface3
				0xeeeeee, 0xf43753, // Overlay0
				0xeeeeee, 0x6a8759, // Overlay1
				0x555555, 0x2b2b2b, // Subtext0
				0x555555, 0x2b2b2b, // Subtext1
				0xf43753, 0x2b2b2b, // Subtext2
			},

			"Nord": {
				0xd8dee9, 0x2e3440, // Terminal
				0xeceff4, 0x3b4252, // Surface0
				0xeceff4, 0x4c566a, // Surface1
				0xeceff4, 0x3b4252, // Surface2
				0xeceff4, 0x5e81ac, // Surface3
				0x2e3440, 0xbf616a, // Overlay0
				0x2e3440, 0xa3be8c, // Overlay1
				0x4c566a, 0x2e3440, // Subtext0
				0x3b4252, 0x2e3440, // Subtext1
				0x5e81ac, 0x2e3440, // Subtext2
			},

			"Matrix": {
				0x008f11, 0x0d0208, // Terminal
				0x003b00, 0x0d0208, // Surface0
				0x00ff41, 0x0d0208, // Surface1
				0x00ff41, 0x0d0208, // Surface2
				0x0d0208, 0x00ff41, // Surface3
				0x0d0208, 0x00ff41, // Overlay0
				0x0d0208, 0x00ff41, // Overlay1
				0x003b00, 0x0d0208, // Subtext0
				0x0d0208, 0x0d0208, // Subtext1
				0x00ff41, 0x0d0208, // Subtext2
			},

			"Monochrome": {
				0xffffff, 0x000000, // Terminal
				0xffffff, 0x000000, // Surface0
				0xffffff, 0x000000, // Surface1
				0xffffff, 0x000000, // Surface2
				0x000000, 0xffffff, // Surface3
				0x000000, 0xffffff, // Overlay0
				0x000000, 0xffffff, // Overlay1
				0xffffff, 0x000000, // Subtext0
				0x000000, 0x000000, // Subtext1
				0x000000, 0xffffff, // Subtext2
			},
		},
		names: []string{
			"Examiner-Dark",
			"Examiner-Light",
			"Catppuccin-Latte",
			"Catppuccin-Frappe",
			"Catppuccin-Macchiato",
			"Catppuccin-Mocha",
			"Solarized-Dark",
			"Solarized-Light",
			"VSCode-Dark",
			"VSCode-Light",
			"Monokai",
			"Darcula",
			"Nord",
			"Matrix",
			"Monochrome",
		},
		index: 0,
	}

	ts := themes.New()

	if ts != nil {
		for _, tt := range ts.Themes {
			t.names = append(t.names, tt.Name)

			t.palettes[tt.Name] = palette{
				tt.Terminal[0], tt.Terminal[1],
				tt.Surface0[0], tt.Surface0[1],
				tt.Surface1[0], tt.Surface1[1],
				tt.Surface2[0], tt.Surface2[1],
				tt.Surface3[0], tt.Surface3[1],
				tt.Overlay0[0], tt.Overlay0[1],
				tt.Overlay1[0], tt.Overlay1[1],
				tt.Subtext0[0], tt.Subtext0[1],
				tt.Subtext1[0], tt.Subtext1[1],
				tt.Subtext2[0], tt.Subtext2[1],
			}
		}
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

	if len(name) == 0 {
		name = Default
	}

	for i, n := range t.names {
		if strings.ToLower(n) == strings.ToLower(name) {
			t.index = i
			break
		}
	}

	if t.index == -1 {
		sys.Error("theme not found")

		t.index = 0
	}

	p := t.palettes[t.names[t.index]]

	Cursor = tcell.NewHexColor(p[4])

	Terminal = newStyle(p[0], p[1])
	Surface0 = newStyle(p[2], p[3])
	Surface1 = newStyle(p[4], p[5])
	Surface2 = newStyle(p[6], p[7])
	Surface3 = newStyle(p[8], p[9])
	Overlay0 = newStyle(p[10], p[11])
	Overlay1 = newStyle(p[12], p[13])
	Subtext0 = newStyle(p[14], p[15])
	Subtext1 = newStyle(p[16], p[17])
	Subtext2 = newStyle(p[18], p[19])
}

func newStyle(fg, bg int32) tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.NewHexColor(fg)).
		Background(tcell.NewHexColor(bg))
}
