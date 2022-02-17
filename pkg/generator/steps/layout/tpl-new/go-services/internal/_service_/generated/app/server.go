package app

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type serverModel struct {
	TplHeader              string
	CoreImportPath         string
	ConfigsImportPath      string
	ApiEndpointEnv         string
	ApiEndpoint            string
	MaintenanceEndpointEnv string
	MaintenanceEndpoint    string
}

func newServerModel(ctx *gencontext.GenContext) (serverModel, error) {
	resolved, err := ctx.EndpointsResolver.ResolveEndpoints(ctx.GetServiceName())
	if err != nil {
		return serverModel{}, err
	}

	// TODO: move paths to description
	return serverModel{
		TplHeader:              ctx.GetWorkspace().TplHeader,
		CoreImportPath:         ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
		ConfigsImportPath:      ctx.GetWorkspace().GetGoConfigsImportPath(),
		ApiEndpointEnv:         endpoints.MakeApiEndpointEnvName(ctx.GetServiceName()),
		ApiEndpoint:            resolved.Api,
		MaintenanceEndpointEnv: endpoints.MakeMaintenanceEndpointEnvName(ctx.GetServiceName()),
		MaintenanceEndpoint:    resolved.Maintenance,
	}, nil
}
