package tpl

import (
	"fmt"
	"strings"
	"unicode"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
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

func (c ServiceModel) GetEndpointEnvName() string {
	return MakeServerEnvName(c.ServiceName)
}

func SanitizeServiceName(serviceName string) string {
	if unicode.IsDigit(rune(serviceName[0])) {
		serviceName = "service_" + serviceName
	}
	serviceName = strings.ReplaceAll(serviceName, "-", "_")

	return serviceName
}

func MakeServerEnvName(serviceName string) string {
	sanitizedName := SanitizeServiceName(serviceName)
	return strings.ToUpper(sanitizedName) + "_SERVER_ENDPOINT"
}
