package openapi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/mifyconfig"
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

func (g *OpenAPIGenerator) doGenerateClient(ctx *gencontext.GenContext, assetsPath string, clientName string, schemaPath string, targetPath string) error {
	generatedPath := filepath.Join(g.basePath, targetPath, "generated", "api", "clients", clientName)

	packageName := MakePackageName(clientName)

	endpoints, err := ctx.EndpointsResolver.ResolveEndpoints(clientName)
	if err != nil {
		return err
	}

	err = runOpenapiGenerator(ctx, g.basePath, schemaPath, assetsPath,
		generatedPath, packageName, clientName, endpoints.Api, g.info)
	if err != nil {
		return fmt.Errorf("failed to run openapi-generator: %w", err)
	}

	// TODO: go specific
	if g.language == mifyconfig.ServiceLanguageGo {
		err = os.Remove(filepath.Join(generatedPath, "api"))
		if err != nil {
			return err
		}
	}

	err = formatGenerated(generatedPath, g.language)
	if err != nil {
		return err
	}

	return nil
}

func (g *OpenAPIGenerator) doRemoveClient(ctx *gencontext.GenContext, clientName string, targetPath string) error {
	generatedPath := filepath.Join(g.basePath, targetPath, "generated", "api", "clients", clientName)

	if err := os.RemoveAll(generatedPath); err != nil {
		return fmt.Errorf("failed to remove client: %w", err)
	}

	return nil
}

func SanitizeServiceName(serviceName string) string {
	if unicode.IsDigit(rune(serviceName[0])) {
		serviceName = "service_" + serviceName
	}
	serviceName = strings.ReplaceAll(serviceName, "-", "_")

	return serviceName
}

func MakePackageName(clientName string) string {
	packageName := SanitizeServiceName(clientName)
	return packageName + "_client"
}

func MakeClientEnvName(serviceName string) string {
	sanitizedName := SanitizeServiceName(serviceName)
	return strings.ToUpper(sanitizedName) + "_CLIENT_ENDPOINT"
}

func SnakeCaseToCamelCase(inputUnderScoreStr string, capitalize bool) (camelCase string) {
	isToUpper := false
	for k, v := range inputUnderScoreStr {
		if k == 0 && capitalize {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}
