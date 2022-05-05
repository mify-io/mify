package openapi

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type appModel struct {
	TplHeader          string
	ServiceName        string
}

func newAppModel(ctx *gencontext.GenContext) appModel {
	return appModel{
		TplHeader:          ctx.GetWorkspace().TplHeaderPy,
		ServiceName:        ctx.GetMifySchema().ServiceName,
	}
}
