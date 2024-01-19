package apputil

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type appUtilModel struct {
	TplHeader      string
	ServiceName    string
	AppImportPath  string
	CoreImportPath string
}

func newAppUtilModel(ctx *gencontext.GenContext) appUtilModel {
	return appUtilModel{
		TplHeader:      ctx.GetWorkspace().TplHeader,
		ServiceName:    ctx.GetServiceName(),
		CoreImportPath: ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
		AppImportPath: fmt.Sprintf(
			"%s/internal/%s/app",
			ctx.GetWorkspace().GetGoModule(),
			ctx.MustGetMifySchema().ServiceName),
	}
}
