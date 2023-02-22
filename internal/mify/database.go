package mify

import (
	"fmt"

	"github.com/mify-io/mify/pkg/workspace/mutators/database"
)

func AddPostgres(ctx *CliContext, basePath string, serviceName string) error {
	err := database.AddPostgres(ctx.MustGetMutatorContext(), serviceName)
	if err != nil {
		return fmt.Errorf("can't add postgres: %w", err)
	}

	if err = ServiceGenerate(ctx, basePath, serviceName, false, false); err != nil {
		return fmt.Errorf("error during generation: %w", err)
	}

	return nil
}

func RemovePostgres(ctx *CliContext, basePath string, serviceName string) error {
	err := database.RemovePostgres(ctx.MustGetMutatorContext(), serviceName)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, basePath, serviceName, false, false)
}
