package openapi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/openapi/processors"
)

func (g *OpenAPIGenerator) doGenerateServer(
	ctx *gencontext.GenContext,
	schemaPath string,
	paths []string,
) error {
	postProcessor, err := processors.NewPostProcessor(g.language)
	if err != nil {
		return err
	}

	generatorConf, err := postProcessor.GetServerGeneratorConfig(ctx)
	if err != nil {
		return err
	}

	endpoints, err := ctx.EndpointsResolver.ResolveEndpoints(ctx.GetServiceName())
	if err != nil {
		return err
	}

	err = runOpenapiGenerator(
		ctx, g.basePath, schemaPath, g.serverAssetsPath,
		generatorConf.TargetPath, generatorConf.PackageName,
		g.info.ServiceName, endpoints.Api, g.info)
	if err != nil {
		return fmt.Errorf("failed to run openapi-generator: %w", err)
	}

	err = postProcessor.ProcessServer(ctx)
	if err != nil {
		return err
	}

	err = postProcessor.PopulateServerHandlers(ctx, paths)
	if err != nil {
		return err
	}

	err = postProcessor.Format(ctx)
	if err != nil {
		return err
	}

	return nil
}

type schemaYaml struct {
	origin map[string]interface{}

	// Parsed fields
	paths map[interface{}]interface{}
}

func (g *OpenAPIGenerator) loadSchema(ctx *gencontext.GenContext, schemaPath string) (schemaYaml, error) {
	schema, err := g.readSchema(ctx, schemaPath)
	if err != nil {
		return schemaYaml{}, fmt.Errorf("failed to read schema: %s: %w", schemaPath, err)
	}

	res, err := parseSchema(schema)
	if err != nil {
		return schemaYaml{}, fmt.Errorf("failed to parse schema '%s': %s", schemaPath, err)
	}

	return res, nil
}

func parseSchema(schema map[string]interface{}) (schemaYaml, error) {
	res := schemaYaml{origin: schema}
	paths, ok := schema["paths"]
	if !ok {
		return res, fmt.Errorf("missing 'paths' section")
	}
	res.paths = paths.(map[interface{}]interface{})

	return res, nil
}

func mergeSchemas(main schemaYaml, generated schemaYaml) (schemaYaml, error) {
	newPaths := make(map[interface{}]interface{})

	for k, v := range main.paths {
		newPaths[k] = v
	}

	for k, v := range generated.paths {
		if _, ok := newPaths[k]; ok {
			return schemaYaml{}, fmt.Errorf("key %s is alredy defined in generated schema (section 'paths'). Please remove it from your schema", k)
		}

		newPaths[k] = v
	}

	newOrigin := main.origin
	newOrigin["paths"] = newPaths

	return parseSchema(newOrigin)
}

func (g *OpenAPIGenerator) makeServerEnrichedSchema(ctx *gencontext.GenContext, schemaDir string) (string, []string, error) {
	mainSchemaPath := filepath.Join(g.basePath, schemaDir, "/api.yaml")
	generatedSchemaPath := filepath.Join(g.basePath, schemaDir, filepath.Join("/", GENERATED_API_FILENAME))

	mainSchema, err := g.loadSchema(ctx, mainSchemaPath)
	if err != nil {
		return "", nil, err
	}

	mergedSchema := mainSchema
	if _, err := os.Stat(generatedSchemaPath); err == nil {
		generatedSchema, err := g.loadSchema(ctx, generatedSchemaPath)
		if err != nil {
			return "", nil, err
		}

		mergedSchema, err = mergeSchemas(mainSchema, generatedSchema)
		if err != nil {
			return "", nil, err
		}
	}

	// TODO mapstructure
	pathsList := []string{}
	for path, v := range mergedSchema.paths {
		ctx.Logger.Infof("processing path: %s", path)
		methods := v.(map[interface{}]interface{})
		if _, ok := methods["$ref"]; ok {
			return "", nil, fmt.Errorf("paths with $ref are not supported yet")
		}
		pathsList = append(pathsList, path.(string))
		for m, vv := range methods {
			ctx.Logger.Infof("processing method: %s", m)
			method := vv.(map[interface{}]interface{})
			pathStr := path.(string)
			pathStr = strings.ReplaceAll(pathStr, "{", "")
			pathStr = strings.ReplaceAll(pathStr, "}", "")
			method["tags"] = []string{pathStr}
			methods[m] = method
		}
	}

	path, err := g.saveEnrichedSchema(ctx, mergedSchema.origin, mainSchemaPath, CACHE_SERVER_SUBDIR)
	if err != nil {
		return "", nil, err
	}
	return path, pathsList, nil
}
