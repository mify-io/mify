package service

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type nuxtConfigModel struct {
	ServiceName string
}

func newNuxtConfigModel(ctx *gencontext.GenContext) nuxtConfigModel {
	return nuxtConfigModel{
		ServiceName: ctx.GetServiceName(),
	}
}
