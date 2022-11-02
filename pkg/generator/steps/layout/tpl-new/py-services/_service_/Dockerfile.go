package service

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type dockerfileModel struct {
	ServiceName    string
	ApiEndpointEnv string
}

func newDockerfileModel(ctx *gencontext.GenContext) dockerfileModel {
	return dockerfileModel{
		ServiceName:    ctx.MustGetMifySchema().ServiceName,
		ApiEndpointEnv: endpoints.MakeApiEndpointEnvName(ctx.GetServiceName()),
	}
}
