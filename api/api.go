package api

import (
	_ "embed"
)

var (
	//go:embed evidence.schema.json
	SchemaJson string

	//go:embed evidence.schema.sql
	SchemaSql string

	//go:embed evidence.schema.xsd
	SchemaXdd string
)
