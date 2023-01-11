package jsservices

import (
	"sort"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

type packageJsonModel struct {
	ServiceList []string
}

func newPackageJsonModel(ctx *gencontext.GenContext) packageJsonModel {
	return packageJsonModel{
		ServiceList: getServiceList(ctx),
	}
}

func getServiceList(ctx *gencontext.GenContext) []string {
	schemas := ctx.GetSchemaCtx().GetAllSchemas()
	res := make([]string, 0)
	for serviceName, schemas := range schemas {
		if schemas.GetMify().Language == mifyconfig.ServiceLanguageJs {
			res = append(res, serviceName)
		}
	}
	sort.Strings(res)
	return res
}
