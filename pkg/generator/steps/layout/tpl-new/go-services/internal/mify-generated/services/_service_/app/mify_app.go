package app

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type mifyAppModel struct {
	TplHeader           string
	ServiceName         string
	AppImportPath       string
	AppRouterImportPath string
	InitImportPath      string
	CoreImportPath      string
	ApiImportPath       string
}

func newMifyAppModel(ctx *gencontext.GenContext) mifyAppModel {
	// TODO: move paths to description
	return mifyAppModel{
		TplHeader:      ctx.GetWorkspace().TplHeader,
		ServiceName:    ctx.GetServiceName(),
		CoreImportPath: ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
		AppImportPath: fmt.Sprintf(
			"%s/internal/%s/app",
			ctx.GetWorkspace().GetGoModule(),
			ctx.MustGetMifySchema().ServiceName),
		AppRouterImportPath: fmt.Sprintf(
			"%s/internal/%s/app/router",
			ctx.GetWorkspace().GetGoModule(),
			ctx.MustGetMifySchema().ServiceName),
		InitImportPath: fmt.Sprintf(
			"%s/api/init",
			ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePackage(),
		),
		ApiImportPath: fmt.Sprintf(
			"%s/api",
			ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePackage(),
		),
	}
}
