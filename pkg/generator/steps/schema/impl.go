package schema

import (
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/schema/context"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace"
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

	files, err := os.ReadDir(ctx.GetWorkspace().GetSchemasRootAbsPath())
	if err != nil {
		return nil, fmt.Errorf("can't iterate schemas directory: %w", err)
	}

	for _, f := range files {
		serviceName := f.Name()
		if f.Name() == workspace.ExternalSchemasDir {
			continue
		}

		openapiSchemas, err := extractOpenapiSchemas(ctx, serviceName, false)
		if err != nil {
			return nil, fmt.Errorf("can't collect openapi schemas for service '%s': %w", serviceName, err)
		}

		mifySchema, err := extractMifySchema(ctx, serviceName)
		if err != nil {
			return nil, fmt.Errorf("can't extract mify schema for service '%s': %w", serviceName, err)
		}

		schemas[serviceName] = context.NewServiceSchemas(openapiSchemas, mifySchema)
	}

	if _, err := os.Stat(ctx.GetWorkspace().GetExternalSchemasRootAbsPath()); os.IsNotExist(err) {
		return schemas, nil
	}

	externalFiles, err := os.ReadDir(ctx.GetWorkspace().GetExternalSchemasRootAbsPath())
	if err != nil {
		return nil, fmt.Errorf("can't iterate external schemas directory: %w", err)
	}
	for _, f := range externalFiles {
		serviceName := f.Name()
		openapiSchemas, err := extractOpenapiSchemas(ctx, serviceName, true)
		if err != nil {
			return nil, fmt.Errorf("can't collect openapi schemas for external service '%s': %w", serviceName, err)
		}

		mifySchema, err := makeExternalServiceMifySchema(ctx, serviceName)
		if err != nil {
			return nil, fmt.Errorf("can't create mify schema for external service '%s': %w", serviceName, err)
		}

		schemas[serviceName] = context.NewServiceSchemas(openapiSchemas, mifySchema)
	}

	return schemas, nil
}

func extractOpenapiSchemas(ctx *gencontext.GenContext, forService string, isExternal bool) (context.OpenapiServiceSchemas, error) {
	openapiSchemasDir := ctx.GetWorkspace().GetApiSchemaDirAbsPath(forService)
	if isExternal {
		openapiSchemasDir = ctx.GetWorkspace().GetExternalApiSchemaDirAbsPath(forService)
	}

	files, err := os.ReadDir(openapiSchemasDir)
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

func makeDefaultDatabaseName(serviceName string) string {
	if unicode.IsDigit(rune(serviceName[0])) {
		serviceName = "db_" + serviceName
	}
	return strings.ToLower(strings.ReplaceAll(serviceName, "-", "_"))
}

func extractMifySchema(ctx *gencontext.GenContext, forService string) (*mifyconfig.ServiceConfig, error) {
	path := ctx.GetWorkspace().GetMifySchemaAbsPath(forService)
	config, err := mifyconfig.ReadServiceCfg(path)
	if err != nil {
		return nil, err
	}
	if len(config.Postgres.DatabaseName) == 0 {
		config.Postgres.DatabaseName = makeDefaultDatabaseName(config.ServiceName)
	}

	return config, nil
}

func makeExternalServiceMifySchema(ctx *gencontext.GenContext, forService string) (*mifyconfig.ServiceConfig, error) {
	return &mifyconfig.ServiceConfig{
		ServiceName: forService,
		IsExternal: true,
	}, nil
}
