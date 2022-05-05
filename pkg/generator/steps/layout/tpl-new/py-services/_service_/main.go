package service

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type mainModel struct {
	ServiceName string
}

func newMainModel(ctx *gencontext.GenContext) mainModel {
	return mainModel{
		ServiceName: ctx.MustGetMifySchema().ServiceName,
	}
}
