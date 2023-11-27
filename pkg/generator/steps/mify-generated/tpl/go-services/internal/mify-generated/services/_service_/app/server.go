package app

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type serverModel struct {
	TplHeader              string
	CoreImportPath         string
	ConfigsImportPath      string
}

func newServerModel(ctx *gencontext.GenContext) serverModel {
	return serverModel{
		TplHeader:              ctx.GetWorkspace().TplHeader,
		CoreImportPath:         ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
		ConfigsImportPath:      ctx.GetWorkspace().GetGoConfigsImportPath(),
	}
}
