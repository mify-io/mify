package processors

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
)

type goPostProcessor struct{}

func newGoProcessor() *goPostProcessor {
	return &goPostProcessor{}
}

func (p *goPostProcessor) GetServerGeneratorConfig(ctx *gencontext.GenContext) (GeneratorConfig, error) {
	basePath := ctx.GetWorkspace().BasePath
	targetPath, err := ctx.GetWorkspace().GetServiceGeneratedAPIRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language)
	if err != nil {
		return GeneratorConfig{}, err
	}
	generatedPath := filepath.Join(basePath, targetPath, "generated")
	return GeneratorConfig{
		TargetPath:  generatedPath,
		PackageName: SERVER_PACKAGE_NAME,
	}, nil
}

func (p *goPostProcessor) GetClientGeneratorConfig(ctx *gencontext.GenContext, clientName string) (GeneratorConfig, error) {
	basePath := ctx.GetWorkspace().BasePath
	targetPath, err := ctx.GetWorkspace().GetServiceGeneratedAPIRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language)
	if err != nil {
		return GeneratorConfig{}, err
	}
	generatedPath := filepath.Join(basePath, targetPath, "generated", "api", "clients", clientName)
	packageName := endpoints.SanitizeServiceName(clientName) + "_client"
	return GeneratorConfig{
		TargetPath:  generatedPath,
		PackageName: packageName,
	}, nil
}

func (p *goPostProcessor) ProcessServer(ctx *gencontext.GenContext) error {
	return nil
}

func (p *goPostProcessor) ProcessClient(ctx *gencontext.GenContext, clientName string) error {
	targetPath, err := ctx.GetWorkspace().GetServiceGeneratedAPIRelPath(ctx.GetServiceName(), ctx.MustGetMifySchema().Language)
	if err != nil {
		return err
	}
	generatedPath := filepath.Join(ctx.GetWorkspace().BasePath, targetPath, "generated", "api", "clients", clientName)
	err = os.Remove(filepath.Join(generatedPath, "api"))
	if err != nil {
		return err
	}
	return nil
}

func (p *goPostProcessor) PopulateServerHandlers(ctx *gencontext.GenContext, paths []string) error {
	var handlerGlob, handlerSuffix, targetServiceName string
	handlerGlob = "api_*_service.go"
	handlerSuffix = "_service.go"
	targetServiceName = "service.go"
	targetDir, err := ctx.GetWorkspace().GetServiceGeneratedAPIRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language)
	if err != nil {
		return err
	}
	generatedPath := filepath.Join(ctx.GetWorkspace().BasePath, targetDir, "generated")
	apiPath := filepath.Join(generatedPath, "api")
	handlersPath := filepath.Join(ctx.GetWorkspace().BasePath, targetDir, "handlers")
	err = sanitizeServerHandlersImports(ctx, apiPath)
	if err != nil {
		return err
	}
	services, err := filepath.Glob(filepath.Join(apiPath, handlerGlob))
	if err != nil {
		return err
	}
	ctx.Logger.Infof("services: %v", services)
	if len(services) == 0 {
		ctx.Logger.Infof("no handlers to move")
		return nil
	}
	pathsSet := map[string]string{}
	for _, path := range paths {
		ctx.Logger.Infof("pre path: %s", path)
		pathsSet[toAPIFilename(path)] = path
	}
	ctx.Logger.Infof("paths: %v", pathsSet)
	for _, service := range services {
		serviceFileName := filepath.Base(service)
		serviceFileName = strings.TrimSuffix(serviceFileName, handlerSuffix)
		path, ok := pathsSet[serviceFileName]
		if !ok {
			return fmt.Errorf("failed to find path for service file: %s", serviceFileName)
		}

		ctx.Logger.Infof("processing handler for: %v", path)
		targetFile := filepath.Join(handlersPath, path, targetServiceName)
		defer func(svc string) {
			if err := os.Remove(svc); err != nil {
				ctx.Logger.Infof("failed to remove service file: %s: %s", svc, err)
				return
			}
			ctx.Logger.Infof("cleaned generated service file: %s", svc)
		}(service)

		if _, err := os.Stat(targetFile); err == nil {
			ctx.Logger.Infof("skipping existing handler for: %v", path)
			continue
		}
		if err := os.MkdirAll(filepath.Join(handlersPath, path), 0755); err != nil {
			return err
		}
		if err := createServerHandlersFile(ctx, service, targetFile); err != nil {
			return err
		}
		ctx.Logger.Infof("created handler for: %v", path)
	}
	return nil
}

func (p *goPostProcessor) Format(ctx *gencontext.GenContext) error {
	targetDir, err := ctx.GetWorkspace().GetServiceGeneratedAPIRelPath(
		ctx.GetServiceName(), ctx.MustGetMifySchema().Language)
	if err != nil {
		return err
	}
	generatedPath := filepath.Join(ctx.GetWorkspace().BasePath, targetDir, "generated")
	apiPath := filepath.Join(generatedPath, "api")
	return filepath.WalkDir(apiPath, func(path string, d fs.DirEntry, ferr error) error {
		if d == nil {
			return fmt.Errorf("failed to format: %s: %w", apiPath, ferr)
		}
		if d.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".go" {
			return nil
		}
		f, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			return fmt.Errorf("failed to format %s: %w", path, err)
		}

		data, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("failed to format %s: %w", path, err)
		}

		fmtData, err := format.Source(data)
		if err != nil {
			return fmt.Errorf("failed to format %s: %w", path, err)
		}

		err = f.Truncate(0)
		if err != nil {
			return fmt.Errorf("failed to format %s: %w", path, err)
		}

		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to format %s: %w", path, err)
		}

		w := bufio.NewWriter(f)
		_, err = w.Write(fmtData)
		if err != nil {
			return fmt.Errorf("failed to format %s: %w", path, err)
		}

		err = w.Flush()
		if err != nil {
			return fmt.Errorf("failed to format %s: %w", path, err)
		}

		return nil
	})
}

// taken from openapi-generator
func isReservedFilename(name string) bool {
	parts := strings.Split(name, "_")
	suffix := parts[len(parts)-1]

	reservedSuffixes := []string{
		// Test
		"test",
		// $GOOS
		"aix", "android", "darwin", "dragonfly", "freebsd", "illumos", "js", "linux", "netbsd", "openbsd",
		"plan9", "solaris", "windows",
		// $GOARCH
		"386", "amd64", "arm", "arm64", "mips", "mips64", "mips64le", "mipsle", "ppc64", "ppc64le", "s390x",
		"wasm",
	}
	reservedSuffixesSet := map[string]struct{}{}
	for _, suf := range reservedSuffixes {
		reservedSuffixesSet[suf] = struct{}{}
	}
	_, ok := reservedSuffixesSet[suffix]
	return ok
}

var (
	pkgSeparatorPattern  = regexp.MustCompile(`\.`)
	dollarPattern        = regexp.MustCompile(`\$`)
)

// taken from openapi-generator
func underscore(word string) string {
	// Replace package separator with slash.
	result := pkgSeparatorPattern.ReplaceAllString(word, "_")
	// Replace $ with two underscores for inner classes.
	result = dollarPattern.ReplaceAllString(result, "__")
	result = endpoints.CamelCaseToSnakeCase(result)
	result = strings.ReplaceAll(result, "-", "_")
	// replace space with underscore
	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ToLower(result)
	return result
}

// taken from openapi-generator
func toAPIFilename(name string) string {
	// NOTE: openapi-generator transforms tag to camelCase, we don't do that here
	// we just remove slashes from path and then use openapi-generator logic
	// to convert this path to filename.
	api := strings.ReplaceAll(name, "{", "")
	api = strings.ReplaceAll(api, "}", "")
	api = strings.TrimPrefix(api, "/")
	api = strings.TrimSuffix(api, "/")
	api = strings.ReplaceAll(api, "/", "_")
	// replace - with _ e.g. created-at => created_at
	api = strings.ReplaceAll(api, "-", "_")
	// // e.g. PetApi.go => pet_api.go
	api = "api_" + underscore(api)
	if isReservedFilename(api) {
		api += "_"
	}
	return api
}

func sanitizeServerHandlersImports(ctx *gencontext.GenContext, apiPath string) error {
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
		if _, err := w.WriteString(line); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	ctx.Logger.Infof("sanitized routes imports")
	return nil
}

func createServerHandlersFile(ctx *gencontext.GenContext, serviceFile string, targetFile string) error {
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

	data, err := os.ReadFile(serviceFile)
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
			if _, err := w.WriteString(line); err != nil {
				return err
			}

			if i+1 < len(lines) && lines[i+1] != sectionStart {
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
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
		ctx.Logger.Infof("writing param %s %s", decl.Name, decl.Type)
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
