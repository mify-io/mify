package tpl

import (
	"strings"
	"unicode"

	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

type ServiceModel struct {
	ServiceName string
	Repository  string
	Language    mifyconfig.ServiceLanguage
	GoModule    string
	Workspace   WorkspaceModel
	ServiceList []string
}

func NewServiceModel(ctx *gencontext.GenContext) *ServiceModel {
	return &ServiceModel{
		ServiceName: ctx.GetServiceName(),
		Repository:  ctx.GetWorkspace().GetRepository(),
		Language:    ctx.MustGetMifySchema().Language,
		GoModule:    ctx.GetWorkspace().GoRoot,
		Workspace:   *NewWorkspaceModel(ctx),
		ServiceList: getServiceList(ctx), // Js only
	}
}

func getServiceList(ctx *gencontext.GenContext) []string {
	schemas := ctx.GetSchemaCtx().GetAllSchemas()
	res := make([]string, 0)
	for serviceName, schemas := range *schemas {
		if schemas.GetMify().Language == mifyconfig.ServiceLanguageJs {
			res = append(res, MakeServerEnvName(serviceName))
		}
	}
	return res
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
