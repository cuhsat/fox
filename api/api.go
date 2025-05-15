package api

import (
	_ "embed"
)

var (
	//go:embed evidence.schema.sql
	SqlSchema string
)
