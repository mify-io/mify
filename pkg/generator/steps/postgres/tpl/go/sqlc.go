package gotpl

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type sqlcModel struct {
	MigrationsDir string
	QueriesDir string
	OutDir string
}

func NewSqlcModel(ctx *gencontext.GenContext) (sqlcModel, error) {
	dbName := ctx.MustGetMifySchema().Postgres.DatabaseName
	// NOTE: sqlc doesn't support absolute paths, these are relative to sqlc.yaml
	// location, so go-services/internal/<service-name>/generated/postgres
	return sqlcModel{
		MigrationsDir: "../../../../../migrations/"+dbName,
		QueriesDir: "../../../../../sql-queries/"+dbName,
		OutDir: ".",
	}, nil
}
