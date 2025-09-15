package fox

import (
	_ "embed"
)

const (
	Product = "Forensic Examiner"
	Website = "forensic-examiner.eu"
	Author  = "Christian Uhsat"
	Email   = "christian@uhsat.de"
)

var (
	//go:embed res/ascii.txt
	Ascii string

	//go:embed res/help.txt
	Help string

	//go:embed res/prompt.txt
	Prompt string

	//go:embed res/version.txt
	Version string
)
