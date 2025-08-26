package configs

import (
	_ "embed"
)

var (
	//go:embed default.yaml
	Default string

	//go:embed plugins.yaml
	Plugins string

	//go:embed themes.yaml
	Themes string
)
