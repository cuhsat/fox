package app

import (
	_ "embed"
)

const (
	Product = "Forensic Examiner"
	Author  = "Christian Uhsat"
	Email   = "christian@uhsat.de"
)

var (
	//go:embed ascii.txt
	Ascii string

	//go:embed help.txt
	Help string

	//go:embed prompt.txt
	Prompt string

	//go:embed version.txt
	Version string
)
