package tpl

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace"
)

type ServiceModel struct {
	ServiceName string
	Repository  string
	Language    mifyconfig.ServiceLanguage
	GoModule    string // first line inside go.mod
	Workspace   WorkspaceModel
}

func NewServiceModel(ctx *gencontext.GenContext) *ServiceModel {
	return &ServiceModel{
		ServiceName: ctx.GetServiceName(),
		Repository:  ctx.GetWorkspace().GetRepository(),
		Language:    ctx.MustGetMifySchema().Language,
		GoModule:    fmt.Sprintf("%s/%s", ctx.GetWorkspace().GetRepository(), workspace.GoServicesDirName),
		Workspace:   *NewWorkspaceModel(ctx),
	}
}

func (c ServiceModel) GetApiEndpointEnvName() string {
	return endpoints.MakeApiEndpointEnvName(c.ServiceName)
}

func (c ServiceModel) GetMaintenanceApiEndpointEnvName() string {
	return endpoints.MakeMaintenanceEndpointEnvName(c.ServiceName)
}
