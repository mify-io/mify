package mify

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/schema"
	"github.com/mify-io/mify/pkg/workspace"
)

func runCommand(cmdname string, args ...string) error {
	cmd := exec.Command(cmdname, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func installMigrate(ctx *CliContext) (string, error) {
	toolPath := filepath.Join(build.Default.GOPATH, "bin", "dbmate")
	if _, err := os.Stat(toolPath); errors.Is(err, os.ErrNotExist) {
		ctx.Logger.Printf("Installing migrate tool dbmate...")
		err := runCommand(
			"go", "install", "github.com/amacneil/dbmate@v1.13.0",
		)
		if err != nil {
			return "", err
		}
	}
	return toolPath, nil
}

func ToolMigrate(
	ctx *CliContext, basePath string, serviceName string, command string, extraArgs []string) error {
	workspace, err := workspace.InitDescription(basePath)
	if err != nil {
		return err
	}
	genContext, err := gencontext.NewGenContext(ctx.Ctx, serviceName, workspace, false)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	schemaStep := schema.NewSchemaStep()
	_, err = schemaStep.Execute(genContext)
	if err != nil {
		return fmt.Errorf("failed to get service schemas: %w", err)
	}

	if !genContext.GetMifySchema().Postgres.Enabled {
		return fmt.Errorf("postgres is not enabled for this service")
	}

	migrateBinPath, err := installMigrate(ctx)
	if err != nil {
		return fmt.Errorf("failed to install migrate command: %w", err)
	}
	dbName := genContext.GetMifySchema().Postgres.DatabaseName
	migrationsDir, err := workspace.GetMigrationsDirectory(dbName, genContext.GetMifySchema().Language)
	if err != nil {
		return fmt.Errorf("failed to get migrations directory: %w", err)
	}

	if _, err := os.Stat(migrationsDir); errors.Is(err, os.ErrNotExist) && command != "new" {
		ctx.Logger.Printf("Migrations directory is not initialized, please create your first migration:")
		ctx.Logger.Printf("mify tool migrate %s new <migration-name>", serviceName)
		return nil
	}

	args := []string{
		"--url", "postgres://user:passwd@localhost:5432/" + dbName + "?sslmode=disable",
		"--no-dump-schema",
		"--migrations-dir", migrationsDir,
		command,
	}
	args = append(args, extraArgs...)
	err = runCommand(migrateBinPath, args...)
	if err != nil {
		return fmt.Errorf("failed to execute migrate command: %w", err)
	}
	return nil
}

// genContext.GetMifySchema().Postgres.DatabaseName
