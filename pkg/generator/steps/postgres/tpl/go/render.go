package gotpl

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	sqlc "github.com/sqlc-dev/sqlc/pkg/cli"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/postgres/tpl/go/config"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func needRunSqlc(migrationsDir string, queriesDir string) (bool, error) {
	hasQueries := false
	err := filepath.WalkDir(queriesDir, func(s string, d fs.DirEntry, e error) error {
		if filepath.Ext(s) == ".sql" {
			hasQueries = true
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	if !hasQueries {
		return false, nil
	}
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return false, fmt.Errorf("failed to generate sql queries, migrations dir doesn't exists")
	}
	return hasQueries, nil
}

func Render(ctx *gencontext.GenContext) error {
	sqlcModel, err := NewSqlcModel(ctx)
	if err != nil {
		return err
	}
	migrationsDir, err := ctx.GetWorkspace().GetMigrationsDirectory(
		ctx.MustGetMifySchema().Postgres.DatabaseName, ctx.MustGetMifySchema().Language)
	if err != nil {
		return err
	}
	queriesDir, err := ctx.GetWorkspace().GetSqlQueriesDirectory(
		ctx.MustGetMifySchema().Postgres.DatabaseName, ctx.MustGetMifySchema().Language)
	if err != nil {
		return err
	}
	sqlcConfPath := filepath.Join(ctx.GetWorkspace().GetGoPostgresConfigAbsPath(ctx.GetServiceName()), "sqlc.yaml")
	err = render.RenderMany(
		templates,
		render.NewFile(
			ctx,
			sqlcConfPath,
		).SetModel(sqlcModel),
		render.NewFile(
			ctx,
			filepath.Join(queriesDir, "queries.sql.example"),
		),
	)
	if err != nil {
		return err
	}

	isSqlcRunNeeded, err := needRunSqlc(migrationsDir, queriesDir)
	if err != nil {
		return err
	}

	if isSqlcRunNeeded {
		rc := sqlc.Run([]string{"--file", sqlcConfPath, "generate"})
		if rc != 0 {
			return fmt.Errorf("sqlc exited with non-zero code: %d", rc)
		}
	}

	return config.Render(ctx)
}
