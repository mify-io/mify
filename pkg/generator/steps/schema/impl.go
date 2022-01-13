package schema

import (
	"io/ioutil"
	"path"

	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
	"github.com/chebykinn/mify/pkg/generator/steps/schema/context"
	"github.com/chebykinn/mify/pkg/workspace"
	"github.com/getkin/kin-openapi/openapi3"
)

func execute(ctx *gencontext.GenContext) (*context.SchemaContext, error) {
	openapiSchemas, err := collectOpenapiSchemas(ctx.GetWorkspace())
	if err != nil {
		return nil, err
	}

	return context.NewSchemaContext(openapiSchemas), nil
}

func collectOpenapiSchemas(workspace *workspace.Description) (context.OpenapiSchemas, error) {
	openapiSchemas := make(context.OpenapiSchemas)

	for _, goService := range workspace.GoServices {
		apiSchmeasDir := workspace.GetApiSchemaDirAbsPath(goService.Name)
		files, err := ioutil.ReadDir(apiSchmeasDir)
		if err != nil {
			return nil, err
		}

		serviceSchemas := make(context.ServiceSchemas)
		for _, f := range files {
			schemaAbsPath := path.Join(apiSchmeasDir, f.Name())
			doc, err := openapi3.NewLoader().LoadFromFile(schemaAbsPath)
			if err != nil {
				return nil, err
			}

			serviceSchemas[f.Name()] = doc
		}

		openapiSchemas[goService.Name] = serviceSchemas
	}

	return openapiSchemas, nil
}
