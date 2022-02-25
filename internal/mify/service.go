package mify

import (
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace"
	"github.com/mify-io/mify/pkg/workspace/mutators/service"
)

func CreateService(ctx *CliContext, basePath string, language string, name string) error {
	mutCtx := ctx.MustGetMutatorContext()
	err := service.CreateService(mutCtx, mifyconfig.ServiceLanguage(language), name)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, basePath, name)
}

func CreateFrontend(ctx *CliContext, basePath string, template string, name string) error {
	mutCtx := ctx.MustGetMutatorContext()
	err := service.CreateFrontend(mutCtx, template, name)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, basePath, name)
}

func CreateApiGateway(ctx *CliContext) error {
	mutCtx := ctx.MustGetMutatorContext()
	err := service.CreateApiGateway(mutCtx)
	if err != nil {
		return err
	}

	return ServiceGenerate(ctx, ctx.WorkspacePath, workspace.ApiGatewayName)
}
