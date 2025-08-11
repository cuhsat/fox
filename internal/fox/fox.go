package fox

import (
	_ "embed"
)

const (
	Product = "Forensic Examiner"
)

var (
	//go:embed version.txt
	Version string

	//go:embed prompt.txt
	Prompt string

	//go:embed help.txt
	Help string

	//go:embed fox.txt
	Fox string
)
