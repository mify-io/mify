package jsservices

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/helpers/js"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

type PackageJsonModel struct {
	ServiceList []string
}

func NewPackageJsonModel(ctx *gencontext.GenContext) PackageJsonModel {
	return PackageJsonModel{
		ServiceList: getServiceList(ctx),
	}
}

func getServiceList(ctx *gencontext.GenContext) []string {
	schemas := ctx.GetSchemaCtx().GetAllSchemas()
	res := make([]string, 0)
	for serviceName, schemas := range schemas {
		if schemas.GetMify().Language == mifyconfig.ServiceLanguageJs {
			res = append(res, js.MakeServerEnvName(serviceName))
		}
	}
	return res
}
