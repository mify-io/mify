package apigateway

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/getkin/kin-openapi/openapi3"
)

const (
	ExtensionKey = "x-mify-public"
)

// service name -> public apis
type PublicApis map[string]apiPaths
type apiPaths map[string]*openapi3.PathItem

func scanPublicApis(ctx *gencontext.GenContext) PublicApis {
	res := make(PublicApis)

	for _, goService := range ctx.GetWorkspace().GoServices {
		if goService.Name == ApiGatewayName {
			continue
		}

		serviceSchemas := ctx.GetSchemaCtx().MustGetServiceSchemas(goService.Name)
		openapiSchema := serviceSchemas.GetOpenapi().GetMainSchema()
		res[goService.Name] = extractPublicAPI(openapiSchema)
	}

	return res
}

func extractPublicAPI(schema *openapi3.T) apiPaths {
	res := make(apiPaths)

	for path, pathItem := range schema.Paths {
		var pathItemCopy openapi3.PathItem = *pathItem
		for method, operation := range pathItem.Operations() {
			if _, ok := operation.Extensions[ExtensionKey]; !ok {
				pathItemCopy.SetOperation(method, nil)
			}
		}

		if len(pathItemCopy.Operations()) > 0 {
			res[path] = &pathItemCopy
		}
	}

	return res
}
