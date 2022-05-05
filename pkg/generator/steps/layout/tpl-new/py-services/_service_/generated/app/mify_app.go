package app

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type mifyAppModel struct {
	TplHeader      string
	ServiceName    string
	//AppImportPath  string
	//InitImportPath string
	//CoreImportPath string
	//ApiImportPath  string
}

func newMifyAppModel(ctx *gencontext.GenContext) mifyAppModel {
	// TODO: move paths to description
	return mifyAppModel{
		TplHeader:      ctx.GetWorkspace().TplHeaderPy,
		ServiceName:    ctx.GetServiceName(),
		//CoreImportPath: ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
		//AppImportPath: fmt.Sprintf(
		//	"%s/internal/%s/app",
		//	ctx.GetWorkspace().GetGoModule(),
		//	ctx.MustGetMifySchema().ServiceName),
		//InitImportPath: fmt.Sprintf(
		//	"%s/internal/%s/generated/api/init",
		//	ctx.GetWorkspace().GetGoModule(),
		//	ctx.MustGetMifySchema().ServiceName),
		//ApiImportPath: fmt.Sprintf(
		//	"%s/internal/%s/generated/api",
		//	ctx.GetWorkspace().GetGoModule(),
		//	ctx.MustGetMifySchema().ServiceName),
	}
}
