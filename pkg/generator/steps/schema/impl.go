package schema

import (
	generator "github.com/chebykinn/mify/pkg/generator/context"
	"github.com/chebykinn/mify/pkg/generator/steps/schema/context"
	"github.com/chebykinn/mify/pkg/workspace"
	"github.com/getkin/kin-openapi/openapi3"
)

func execute(ctx *generator.GenContext) (*context.SchemaContext, error) {
	openapiSchemas, err := collectOpenapiSchemas(ctx.GetWorkspace())
	if err != nil {
		return nil, err
	}

	return context.NewSchemaContext(openapiSchemas), nil
}

func collectOpenapiSchemas(workspace *workspace.Description) (context.OpenapiSchemas, error) {
	openapiSchemas := make(context.OpenapiSchemas)

	for _, goService := range workspace.GoServices {
		doc, err := openapi3.NewLoader().LoadFromFile(workspace.GetApiSchemaAbsPath(goService.Name))
		if err != nil {
			return nil, err
		}

		openapiSchemas[goService.Name] = doc
	}

	return openapiSchemas, nil
}
