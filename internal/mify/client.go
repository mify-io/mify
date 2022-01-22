package mify

import (
	"github.com/mify-io/mify/pkg/workspace/mutators/client"
)

func AddClient(ctx *CliContext, basePath string, name string, clientName string) error {
	mutCtx, err := initMutatorCtx(ctx, basePath)
	if err != nil {
		return err
	}

	err = client.AddClient(mutCtx, name, clientName)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, basePath, name)
}

func RemoveClient(ctx *CliContext, basePath string, name string, clientName string) error {
	mutCtx, err := initMutatorCtx(ctx, basePath)
	if err != nil {
		return err
	}

	err = client.RemoveClient(mutCtx, name, clientName)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, basePath, name)
}
