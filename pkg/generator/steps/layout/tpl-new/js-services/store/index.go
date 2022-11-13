package store

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type indexModel struct {
	ServiceName string
}

func newIndexModel(ctx *gencontext.GenContext) indexModel {
	return indexModel{
		ServiceName: ctx.GetServiceName(),
	}
}
