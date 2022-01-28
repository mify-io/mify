package service

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type NuxtConfigModel struct {
	ServiceName string
}

func newNuxtConfigModel(ctx *gencontext.GenContext) NuxtConfigModel {
	return NuxtConfigModel{
		ServiceName: ctx.GetServiceName(),
	}
}
