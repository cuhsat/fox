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

	//go:embed usage.txt
	Usage string

	//go:embed help.txt
	Help string
)
