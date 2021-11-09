package generate

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/util"
)

func (g *OpenAPIGenerator) doGenerateServer(ctx *util.JobPoolContext, schemaPath string, targetPath string) error {
	langStr := string(GENERATOR_LANGUAGE_GO)
	path, err := config.DumpAssets(g.basePath, "openapi/server-template/"+langStr, "openapi/server-template")
	if err != nil {
		return fmt.Errorf("failed to dump assets: %w", err)
	}
	ctx.Logger.Printf("dumped path: %s\n", path)

	generatedPath := filepath.Join(g.basePath, targetPath, "generated")

	err = runOpenapiGenerator(ctx, g.basePath, schemaPath, path, generatedPath, SERVER_PACKAGE_NAME, g.info)
	if err != nil {
		return fmt.Errorf("failed to run openapi-generator: %w", err)
	}

	apiPath := filepath.Join(generatedPath, "api")
	err = sanitizeServerHandlersImports(ctx, apiPath)
	if err != nil {
		return err
	}

	handlersPath := filepath.Join(g.basePath, targetPath, "handlers")
	err = moveServerHandlers(ctx, apiPath, handlersPath)
	if err != nil {
		return err
	}

	err = formatGenerated(apiPath)
	if err != nil {
		return err
	}

	return nil
}

// FIXME: go-specific
func moveServerHandlers(ctx *util.JobPoolContext, apiPath string, handlersPath string) error {
	services, err := filepath.Glob(filepath.Join(apiPath, "api_*_service.go"))
	if err != nil {
		return err
	}
	ctx.Logger.Printf("services: %v\n", services)
	if len(services) == 0 {
		ctx.Logger.Printf("no handlers to move\n")
		return nil
	}
	for _, service := range services {
		baseName := filepath.Base(service)
		baseName = strings.TrimPrefix(baseName, "api_")
		baseName = strings.TrimSuffix(baseName, "_service.go")
		baseName = strings.ReplaceAll(baseName, "_", "/")

		ctx.Logger.Printf("processing handler for: %v\n", baseName)
		targetFile := filepath.Join(handlersPath, baseName, "service.go")
		defer func(svc string) {
			if err := os.Remove(svc); err != nil {
				ctx.Logger.Printf("failed to remove service file: %s: %s\n", svc, err)
				return
			}
			ctx.Logger.Printf("cleaned generated service file: %s\n", svc)
		}(service)

		if _, err := os.Stat(targetFile); err == nil {
			ctx.Logger.Printf("skipping existing handler for: %v", baseName)
			continue
		}
		if err := os.MkdirAll(filepath.Join(handlersPath, baseName), 0755); err != nil {
			return err
		}
		if err := createServerHandlersFile(ctx, service, targetFile); err != nil {
			return err
		}
		ctx.Logger.Printf("created handler for: %v\n", baseName)
	}
	return nil
}

func (g *OpenAPIGenerator) makeServerEnrichedSchema(ctx *util.JobPoolContext, schemaPath string) (string, error) {
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
			method["tags"] = []string{path.(string)}
			methods[m] = method
		}
	}

	return g.saveEnrichedSchema(ctx, doc, schemaPath, CACHE_SERVER_SUBDIR)
}

// FIXME: go-specific
func createServerHandlersFile(ctx *util.JobPoolContext, serviceFile string, targetFile string) error {
	const (
		sectionStart = "// service_params_start"
		sectionEnd   = "// service_params_end"
	)
	var (
		primitivesList = map[string]struct{}{
			"string":      {},
			"bool":        {},
			"uint":        {},
			"uint32":      {},
			"uint64":      {},
			"int":         {},
			"int32":       {},
			"int64":       {},
			"float32":     {},
			"float64":     {},
			"complex64":   {},
			"complex128":  {},
			"rune":        {},
			"byte":        {},
			"interface{}": {},
		}
	)

	data, err := ioutil.ReadFile(serviceFile)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// FIXME: instead of using mustache to generate argument list
	// they're encoded as json so we can process each one individually and
	// add package prefix to non-primitive types.
	var decl struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	buf := bytes.NewBufferString("")
	w := bufio.NewWriter(buf)

	isSectionStart := false
	for i, line := range lines {
		if line == sectionStart {
			isSectionStart = true
			continue
		}
		if line == sectionEnd {
			isSectionStart = false
			continue
		}
		if !isSectionStart {
			w.WriteString(line)
			if i+1 < len(lines) && lines[i+1] != sectionStart {
				w.WriteByte('\n')
			}
			continue
		}
		err := json.Unmarshal([]byte(line), &decl)
		if err != nil {
			return err
		}
		innermostType := decl.Type
		innermostType = strings.ReplaceAll(innermostType, "map[string]", "")
		innermostType = strings.ReplaceAll(innermostType, "[]", "")
		if _, ok := primitivesList[innermostType]; !ok {
			decl.Type = strings.ReplaceAll(decl.Type, innermostType, "openapi."+innermostType)
		}
		ctx.Logger.Printf("writing param %s %s\n", decl.Name, decl.Type)
		fmt.Fprintf(w, ", %s %s", decl.Name, decl.Type)
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	err = os.WriteFile(targetFile, out, 0666)
	if err != nil {
		return err
	}

	return nil
}

// FIXME: go-specific
func sanitizeServerHandlersImports(ctx *util.JobPoolContext, apiPath string) error {
	routesFilePath := filepath.Join(apiPath, "init/routes.go")
	if _, err := os.Stat(routesFilePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("routes file doesn't exists: %s: %w", routesFilePath, err)
	}

	f, err := os.OpenFile(routesFilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	isImportStart := false
	for i, line := range lines {
		if strings.HasPrefix(line, "import") {
			isImportStart = true
			continue
		}
		if isImportStart && strings.HasPrefix(line, ")") {
			break
		}
		if !isImportStart {
			continue
		}
		lines[i] = strings.ReplaceAll(lines[i], "{", "")
		lines[i] = strings.ReplaceAll(lines[i], "}", "")
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	for _, line := range lines {
		w.WriteString(line)
		w.WriteByte('\n')
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	ctx.Logger.Printf("sanitized routes imports\n")
	return nil
}

