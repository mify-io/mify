package schema

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/getkin/kin-openapi/openapi3"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/schema/context"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

func execute(ctx *gencontext.GenContext) (*context.SchemaContext, error) {
	schemas, err := collectSchemas(ctx)
	if err != nil {
		return nil, err
	}

	return context.NewSchemaContext(schemas), nil
}

func collectSchemas(ctx *gencontext.GenContext) (context.AllSchemas, error) {
	schemas := make(context.AllSchemas)

	files, err := ioutil.ReadDir(ctx.GetWorkspace().GetSchemasRootAbsPath())
	if err != nil {
		return nil, fmt.Errorf("can't iterate schemas directory: %w", err)
	}

	for _, f := range files {
		serviceName := f.Name()

		openapiSchemas, err := extractOpenapiSchemas(ctx, serviceName)
		if err != nil {
			return nil, fmt.Errorf("can't collect openapi schemas for service '%s': %w", serviceName, err)
		}

		mifySchema, err := extractMifySchema(ctx, serviceName)
		if err != nil {
			return nil, fmt.Errorf("can't extract mify schema for service '%s': %w", serviceName, err)
		}

		schemas[serviceName] = context.NewServiceSchemas(openapiSchemas, mifySchema)
	}

	return schemas, nil
}

func extractOpenapiSchemas(ctx *gencontext.GenContext, forService string) (context.OpenapiServiceSchemas, error) {
	openapiSchemasDir := ctx.GetWorkspace().GetApiSchemaDirAbsPath(forService)

	files, err := ioutil.ReadDir(openapiSchemasDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	serviceSchemas := make(context.OpenapiServiceSchemas, len(files))
	for _, f := range files {
		schemaAbsPath := path.Join(openapiSchemasDir, f.Name())
		doc, err := openapi3.NewLoader().LoadFromFile(schemaAbsPath)
		if err != nil {
			return nil, err
		}

		serviceSchemas[f.Name()] = doc
	}

	if len(serviceSchemas) == 0 {
		return nil, fmt.Errorf("openapi directory %s exists, but empty", openapiSchemasDir)
	}

	return serviceSchemas, nil
}

func extractMifySchema(ctx *gencontext.GenContext, forService string) (*mifyconfig.ServiceConfig, error) {
	path := ctx.GetWorkspace().GetMifySchemaAbsPath(forService)
	config, err := mifyconfig.ReadServiceCfg(path)
	if err != nil {
		return nil, err
	}

	return config, nil
}
