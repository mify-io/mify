package openapi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
	"github.com/mify-io/mify/pkg/generator/steps/openapi/processors"
)

func (g *OpenAPIGenerator) makeClientEnrichedSchema(ctx *gencontext.GenContext, schemaPath string) (string, error) {
	doc, err := g.readSchema(ctx, schemaPath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema: %s: %w", schemaPath, err)
	}

	pathsIface, ok := doc["paths"]
	if !ok {
		return "", fmt.Errorf("missing paths in schema: %s", schemaPath)
	}
	// TODO mapstructure
	paths := pathsIface.(map[interface{}]interface{})
	for path, v := range paths {
		ctx.Logger.Infof("processing path: %s", path)
		methods := v.(map[interface{}]interface{})
		if _, ok := methods["$ref"]; ok {
			return "", fmt.Errorf("paths with $ref are not supported yet")
		}
		for m, vv := range methods {
			ctx.Logger.Infof("processing method: %s", m)
			method := vv.(map[interface{}]interface{})
			method["tags"] = []string{"api"}
			methods[m] = method
		}
	}

	return g.saveEnrichedSchema(ctx, doc, schemaPath, CACHE_CLIENT_SUBDIR)
}

func (g *OpenAPIGenerator) doGenerateClient(
	ctx *gencontext.GenContext, clientName string, schemaPath string) error {
	endpoints, err := ctx.EndpointsResolver.ResolveEndpoints(clientName)
	if err != nil {
		return err
	}

	postProcessor, err := processors.NewPostProcessor(g.language)
	if err != nil {
		return err
	}

	generatorConf, err := postProcessor.GetClientGeneratorConfig(ctx, clientName)
	if err != nil {
		return err
	}

	err = runOpenapiGenerator(ctx, g.basePath, schemaPath, g.clientAssetsPath,
		generatorConf.TargetPath, generatorConf.PackageName, clientName, endpoints.Api, g.info)
	if err != nil {
		return fmt.Errorf("failed to run openapi-generator: %w", err)
	}

	err = postProcessor.ProcessClient(ctx, clientName)
	if err != nil {
		return err
	}

	err = postProcessor.Format(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (g *OpenAPIGenerator) doRemoveClient(ctx *gencontext.GenContext, clientName string) error {
	targetPath, err := ctx.GetWorkspace().GetServiceGeneratedAPIRelPath(ctx.GetServiceName(), ctx.MustGetMifySchema().Language)
	if err != nil {
		return err
	}
	generatedPath := filepath.Join(g.basePath, targetPath, "generated", "api", "clients", clientName)

	if err := os.RemoveAll(generatedPath); err != nil {
		return fmt.Errorf("failed to remove client: %w", err)
	}

	return nil
}

func MakeClientEnvName(serviceName string) string {
	sanitizedName := endpoints.SanitizeServiceName(serviceName)
	return strings.ToUpper(sanitizedName) + "_CLIENT_ENDPOINT"
}
