package openapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mify-io/mify/internal/mify/util"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type serviceGenCache struct {
	ListenPort int `yaml:"listen_port"`
}

func (g *OpenAPIGenerator) doGenerateServer(
	ctx *gencontext.GenContext, assetsPath string, schemaPath string, targetPath string, paths []string) error {
	generatedPath := filepath.Join(g.basePath, targetPath, "generated")

	listenPort, err := makeServicePort(ctx.GetWorkspace().GetServiceCacheDirectory(ctx.GetServiceName()))
	if err != nil {
		return fmt.Errorf("failed to get service port: %w", err)
	}

	err = runOpenapiGenerator(ctx, g.basePath, schemaPath, assetsPath,
		generatedPath, SERVER_PACKAGE_NAME, g.info.ServiceName, listenPort, g.info)
	if err != nil {
		return fmt.Errorf("failed to run openapi-generator: %w", err)
	}

	apiPath := filepath.Join(generatedPath, "api")
	err = sanitizeServerHandlersImports(ctx, apiPath)
	if err != nil {
		return err
	}

	handlersPath := filepath.Join(g.basePath, targetPath, "handlers")
	err = moveServerHandlers(ctx, apiPath, handlersPath, paths)
	if err != nil {
		return err
	}

	err = formatGenerated(apiPath, g.language)
	if err != nil {
		return err
	}

	return nil
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
	capitalLetterPattern = regexp.MustCompile(`([A-Z]+)([A-Z][a-z][a-z]+)`)
	lowercasePattern     = regexp.MustCompile(`([a-z\d])([A-Z])`)
	pkgSeparatorPattern  = regexp.MustCompile(`\.`)
	dollarPattern        = regexp.MustCompile(`\$`)
)

// taken from openapi-generator
func underscore(word string) string {
	replacementPattern := "$1_$2"
	// Replace package separator with slash.
	result := pkgSeparatorPattern.ReplaceAllString(word, "/")
	// Replace $ with two underscores for inner classes.
	result = dollarPattern.ReplaceAllString(result, "__")
	// Replace capital letter with _ plus lowercase letter.
	result = capitalLetterPattern.ReplaceAllString(result, replacementPattern)
	result = lowercasePattern.ReplaceAllString(result, replacementPattern)
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
	api := strings.TrimPrefix(name, "/")
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

// FIXME: go-specific
func moveServerHandlers(ctx *gencontext.GenContext, apiPath string, handlersPath string, apiPaths []string) error {
	services, err := filepath.Glob(filepath.Join(apiPath, "api_*_service.go"))
	if err != nil {
		return err
	}
	ctx.Logger.Infof("services: %v", services)
	if len(services) == 0 {
		ctx.Logger.Infof("no handlers to move")
		return nil
	}
	pathsSet := map[string]string{}
	for _, path := range apiPaths {
		path = strings.ReplaceAll(path, "{", "")
		path = strings.ReplaceAll(path, "}", "")
		pathsSet[toAPIFilename(path)] = path
	}
	ctx.Logger.Infof("paths: %v", pathsSet)
	for _, service := range services {
		serviceFileName := filepath.Base(service)
		serviceFileName = strings.TrimSuffix(serviceFileName, "_service.go")
		path, ok := pathsSet[serviceFileName]
		if !ok {
			return fmt.Errorf("failed to find path for service file: %s", serviceFileName)
		}

		ctx.Logger.Infof("processing handler for: %v", path)
		targetFile := filepath.Join(handlersPath, path, "service.go")
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
			method["tags"] = []string{path.(string)}
			methods[m] = method
		}
	}

	path, err := g.saveEnrichedSchema(ctx, mergedSchema.origin, mainSchemaPath, CACHE_SERVER_SUBDIR)
	if err != nil {
		return "", nil, err
	}
	return path, pathsList, nil
}

// FIXME: go-specific
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

// FIXME: go-specific
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

func makeServicePort(tmpDir string) (int, error) {
	cacheFilePath := filepath.Join(tmpDir, ".service-cache.yaml")

	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return 0, fmt.Errorf("failed to create service cache directory: %w", err)
	}

	var cache serviceGenCache
	yd := util.NewYAMLData(cacheFilePath)
	err = yd.ReadFile(&cache)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("failed to read service gen cache: %w", err)
	}
	if err == nil && cache.ListenPort > 0 {
		return cache.ListenPort, nil
	}

	port, err := getFreePort()
	if err != nil {
		return 0, fmt.Errorf("failed to get free port: %w", err)
	}
	cache.ListenPort = port

	err = yd.SaveFile(&cache)
	if err != nil {
		return 0, fmt.Errorf("failed to save service gen cache: %w", err)
	}

	return cache.ListenPort, nil
}

func getServicePort(tmpDir string) (int, error) {
	cacheFilePath := filepath.Join(tmpDir, ".service-cache.yaml")

	var cache serviceGenCache
	yd := util.NewYAMLData(cacheFilePath)
	err := yd.ReadFile(&cache)
	if err != nil {
		return 0, fmt.Errorf("failed to read service gen cache: %w", err)
	}
	return cache.ListenPort, nil
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
