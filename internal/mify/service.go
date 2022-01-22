package mify

import (
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace/mutators/service"
)

func CreateService(ctx *CliContext, basePath string, language string, name string) error {
	mutCtx, err := initMutatorCtx(ctx, basePath)
	if err != nil {
		return err
	}

	err = service.CreateService(mutCtx, mifyconfig.ServiceLanguage(language), name)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, basePath, name)
}

func CreateFrontend(ctx *CliContext, basePath string, template string, name string) error {
	mutCtx, err := initMutatorCtx(ctx, basePath)
	if err != nil {
		return err
	}

	err = service.CreateFrontend(mutCtx, template, name)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, basePath, name)
}
