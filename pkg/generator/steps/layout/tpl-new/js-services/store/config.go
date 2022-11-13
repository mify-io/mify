package store

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type configModel struct {
	ServiceName string
}

func newConfigModel(ctx *gencontext.GenContext) configModel {
	return configModel{
		ServiceName: ctx.GetServiceName(),
	}
}
