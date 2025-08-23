package configs

import (
	_ "embed"
)

var (
	//go:embed fox_history
	History string

	//go:embed fox_plugins
	Plugins string

	//go:embed fox_themes
	Themes string

	//go:embed foxrc
	Config string
)
