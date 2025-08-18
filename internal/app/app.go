package app

import (
	_ "embed"
)

const (
	Company = "Hinterland Forensics"
	Product = "Forensic Examiner"
	Url = "https://hiforensics.eu"
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
