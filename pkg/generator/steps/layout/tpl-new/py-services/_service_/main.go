package service

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type mainModel struct {
	ServiceName string
	ServicePackageName string
}

func newMainModel(ctx *gencontext.GenContext) mainModel {
	serviceName := ctx.MustGetMifySchema().ServiceName
	return mainModel{
		ServiceName: serviceName,
		ServicePackageName: endpoints.SanitizeServiceName(serviceName),
	}
}
