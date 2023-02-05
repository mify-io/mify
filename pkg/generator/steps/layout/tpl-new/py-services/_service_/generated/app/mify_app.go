package app

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type mifyAppModel struct {
	TplHeader   string
	ServiceName string
}

func newMifyAppModel(ctx *gencontext.GenContext) mifyAppModel {
	// TODO: move paths to description
	return mifyAppModel{
		TplHeader:   ctx.GetWorkspace().TplHeaderPy,
		ServiceName: ctx.GetServiceName(),
	}
}
