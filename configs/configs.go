package configs

import _ "embed"

var (
	//go:embed default.toml
	Default string

	//go:embed plugins.toml
	Plugins string

	//go:embed themes.toml
	Themes string
)
