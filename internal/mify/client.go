package mify

import (
	"fmt"

	"github.com/mify-io/mify/pkg/workspace/mutators/client"
)

func AddClient(ctx *CliContext, basePath string, name string, clientName string) error {
	err := client.AddClient(ctx.MustGetMutatorContext(), name, clientName)
	if err != nil {
		return fmt.Errorf("can't add client: %w", err)
	}

	if err = ServiceGenerate(ctx, basePath, name, false); err != nil {
		return fmt.Errorf("error during generation: %w", err)
	}

	return nil
}

func RemoveClient(ctx *CliContext, basePath string, name string, clientName string) error {
	err := client.RemoveClient(ctx.MustGetMutatorContext(), name, clientName)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, basePath, name, false)
}
