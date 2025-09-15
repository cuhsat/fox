package info

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
	//go:embed text/ascii.txt
	Ascii string

	//go:embed text/help.txt
	Help string

	//go:embed text/prompt.txt
	Prompt string

	//go:embed text/version.txt
	Version string
)
