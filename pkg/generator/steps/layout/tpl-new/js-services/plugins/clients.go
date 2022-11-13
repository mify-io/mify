package plugins

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type clientsModel struct {
	ServiceName string
}

func newClientsModel(ctx *gencontext.GenContext) clientsModel {
	return clientsModel{
		ServiceName: ctx.GetServiceName(),
	}
}
