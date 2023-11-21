package render

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type ServiceModel struct {
	Name string
}

func NewServiceModel(ctx *gencontext.GenContext) *ServiceModel {
	return &ServiceModel{
		Name: ctx.GetServiceName(),
	}
}

func (c ServiceModel) GetApiEndpointEnvName() string {
	return endpoints.MakeApiEndpointEnvName(c.Name)
}

func (c ServiceModel) GetMaintenanceApiEndpointEnvName() string {
	return endpoints.MakeMaintenanceEndpointEnvName(c.Name)
}
