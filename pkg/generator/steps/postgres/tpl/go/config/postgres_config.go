package config

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type PostgresConfigModel struct {
	ConfigsImportPath string
	DatabaseName string
}

func NewPostgresConfigModel(ctx *gencontext.GenContext) PostgresConfigModel {
	return PostgresConfigModel{
		ctx.GetWorkspace().GetGoConfigsImportPath(),
		ctx.GetMifySchema().Postgres.DatabaseName,
	}
}
