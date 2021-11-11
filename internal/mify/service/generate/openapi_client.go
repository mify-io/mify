package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/chebykinn/mify/internal/mify/core"
)

func (g *OpenAPIGenerator) makeClientEnrichedSchema(ctx *core.Context, schemaPath string) (string, error) {
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
		ctx.Logger.Printf("processing path: %s\n", path)
		methods := v.(map[interface{}]interface{})
		if _, ok := methods["$ref"]; ok {
			return "", fmt.Errorf("paths with $ref are not supported yet")
		}
		for m, vv := range methods {
			ctx.Logger.Printf("processing method: %s\n", m)
			method := vv.(map[interface{}]interface{})
			method["tags"] = []string{"api"}
			methods[m] = method
		}
	}

	return g.saveEnrichedSchema(ctx, doc, schemaPath, CACHE_CLIENT_SUBDIR)
}

func (g *OpenAPIGenerator) doGenerateClient(ctx *core.Context, assetsPath string, clientName string, schemaPath string, targetPath string) error {
	generatedPath := filepath.Join(g.basePath, targetPath, "generated", "api", "clients", clientName)

	packageName := MakePackageName(clientName)

	err := runOpenapiGenerator(ctx, g.basePath, schemaPath, assetsPath, generatedPath, packageName, g.info)
	if err != nil {
		return fmt.Errorf("failed to run openapi-generator: %w", err)
	}

	err = os.Remove(filepath.Join(generatedPath, "api"))
	if err != nil {
		return err
	}

	err = formatGenerated(generatedPath)
	if err != nil {
		return err
	}

	return nil
}

func SanitizeClientName(clientName string) string {
	if unicode.IsDigit(rune(clientName[0])) {
		clientName = "service_" + clientName
	}
	clientName = strings.ReplaceAll(clientName, "-", "_")

	return clientName
}

func MakePackageName(clientName string) string {
	packageName := SanitizeClientName(clientName)
	return packageName + "_client"
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

func CapitalizeFirst(camelCase string) string {
	out := []rune(camelCase)
	out[0] = unicode.ToUpper(rune(camelCase[0]))
	return string(out)
}
