package render

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

type GoServiceModel struct {
	Name string
}

type WorkspaceModel struct {
	Name       string
	MifyGeneratedCommonPackage string
	MifyGeneratedServicePackage string

	PackageName string
}

func NewWorkspaceModel(context *gencontext.GenContext) *WorkspaceModel {
	mifyGen := context.GetWorkspace().GetMifyGenerated(context.MustGetMifySchema())
	return &WorkspaceModel{
		Name:       context.GetWorkspace().Name,
		MifyGeneratedCommonPackage: mifyGen.GetCommonPackage(),
		MifyGeneratedServicePackage: mifyGen.GetServicePackage(),

		PackageName: getPackageName(context),
	}
}

func getPackageName(ctx *gencontext.GenContext) string {
	switch(ctx.MustGetMifySchema().Language) {
	case mifyconfig.ServiceLanguageGo:
		return getGoPackageName(ctx.GetWorkspace().Config)
	case mifyconfig.ServiceLanguageJs:
		// TODO update
		return ctx.MustGetMifySchema().ServiceName
	case mifyconfig.ServiceLanguagePython:
		// TODO update
		return ctx.MustGetMifySchema().ServiceName
	}
	panic(fmt.Sprintf("unknown language: %s", ctx.MustGetMifySchema().Language))
}

func getGoPackageName(config mifyconfig.WorkspaceConfig) string {
	return fmt.Sprintf("%s/%s/%s/go-services",
		config.GitHost,
		config.GitNamespace,
		config.GitRepository)
}
