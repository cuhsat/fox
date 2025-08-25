package configs

import (
	_ "embed"
)

var (
	//go:embed foxrc
	Default string

	//go:embed history
	History string

	//go:embed plugins
	Plugins string

	//go:embed themes
	Themes string
)
