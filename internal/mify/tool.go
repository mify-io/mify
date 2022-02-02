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
	toolPath := filepath.Join(build.Default.GOPATH, "bin", "migrate")
	if _, err := os.Stat(toolPath); errors.Is(err, os.ErrNotExist) {
		ctx.Logger.Printf("Installing migrate...")
		err := runCommand(
			"go", "install", "-tags", "postgres",
			"github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.1",
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
	genContext := gencontext.NewGenContext(ctx.Ctx, serviceName, workspace)
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

	if command == "create" {
		extraArgs = append([]string{
			"-dir", migrationsDir,
			"-ext", "sql",
		}, extraArgs...)
	}
	if _, err := os.Stat(migrationsDir); errors.Is(err, os.ErrNotExist) && command != "create" {
		ctx.Logger.Printf("Migrations directory is not initialized, please create your first migration:")
		ctx.Logger.Printf("mify tool migrate %s create -- -seq <migration-name>", serviceName)
		return nil
	}

	args := []string{
		"-database", "postgres://user:passwd@localhost:5432/" + dbName + "?sslmode=disable",
		"-path", migrationsDir,
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
