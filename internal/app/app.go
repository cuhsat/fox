package app

import (
	_ "embed"
)

const (
	Product = "Forensic Examiner"
	Author  = "Christian Uhsat"
	Email   = "fox@uhsat.de"
)

var (
	//go:embed ascii.txt
	Ascii string

	//go:embed base.txt
	Base string

	//go:embed help.txt
	Help string

	//go:embed version.txt
	Version string
)
