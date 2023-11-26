package render

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type ServiceModel struct {
	Name string
	ApiEndpointEnv         string
	ApiEndpoint            string
	MaintenanceEndpointEnv string
	MaintenanceEndpoint    string
}

func NewServiceModel(ctx *gencontext.GenContext) *ServiceModel {
	resolved, err := ctx.EndpointsResolver.ResolveEndpoints(ctx.GetServiceName())
	if err != nil {
		ctx.Logger.Panic("failed to resolve service endpoints during template generation")
	}
	return &ServiceModel{
		Name: ctx.GetServiceName(),
		ApiEndpointEnv:         endpoints.MakeApiEndpointEnvName(ctx.GetServiceName()),
		ApiEndpoint:            resolved.Api,
		MaintenanceEndpointEnv: endpoints.MakeMaintenanceEndpointEnvName(ctx.GetServiceName()),
		MaintenanceEndpoint:    resolved.Maintenance,
	}
}

func (c ServiceModel) GetApiEndpointEnvName() string {
	return endpoints.MakeApiEndpointEnvName(c.Name)
}

func (c ServiceModel) GetMaintenanceApiEndpointEnvName() string {
	return endpoints.MakeMaintenanceEndpointEnvName(c.Name)
}
