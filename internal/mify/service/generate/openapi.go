package generate

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"

	"os/user"
	"path/filepath"
	"strings"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/util/docker"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

type GeneratorLanguage string

const (
	GENERATOR_LANGUAGE_GO GeneratorLanguage = "go"
)

type OpenAPIGeneratorInfo struct {
	GitHost       string
	GitNamespace  string
	GitRepository string
	ServiceName   string
}

type OpenAPIGenerator struct {
	basePath  string
	schemaDir string
	language  GeneratorLanguage
	info      OpenAPIGeneratorInfo
}

func NewOpenAPIGenerator(basePath string, schemaDir string, language GeneratorLanguage, info OpenAPIGeneratorInfo) OpenAPIGenerator {
	return OpenAPIGenerator{
		basePath:  basePath,
		schemaDir: schemaDir,
		language:  language,
		info:      info,
	}
}

func (g *OpenAPIGenerator) GenerateServer(outputDir string) error {
	// TODO: maybe pass context from caller
	ctx := context.Background()
	schemaPath, err := g.makeEnrichedSchema(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	err = g.doGenerate(ctx, schemaPath, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	return nil

}

// private

func (g *OpenAPIGenerator) makeEnrichedSchema(ctx context.Context) (string, error) {
	schemaPath := filepath.Join(g.basePath, g.schemaDir, "/api.yaml")

	data, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema: %s: %w", schemaPath, err)
	}
	loader := &openapi3.Loader{
		Context:               ctx,
		IsExternalRefsAllowed: true,
	}
	url, err := url.Parse(schemaPath)
	if err != nil {
		return "", fmt.Errorf("failed to validate schema: %s: %w", schemaPath, err)
	}

	openapiDoc, err := loader.LoadFromDataWithPath(data, url)
	if err != nil {
		return "", fmt.Errorf("failed to validate schema: %s: %w", schemaPath, err)
	}
	err = openapiDoc.Validate(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate schema: %s: %w", schemaPath, err)
	}

	doc := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("failed to parse schema: %s: %w", schemaPath, err)
	}

	pathsIface, ok := doc["paths"]
	if !ok {
		return "", fmt.Errorf("missing paths in schema: %s", schemaPath)
	}
	// TODO mapstructure
	paths := pathsIface.(map[interface{}]interface{})
	for path, v := range paths {
		fmt.Printf("debug: processing path: %s\n", path)
		methods := v.(map[interface{}]interface{})
		if _, ok := methods["$ref"]; ok {
			return "", fmt.Errorf("paths with $ref are not supported yet")
		}
		for m, vv := range methods {
			fmt.Printf("debug: processing method: %s\n", m)
			method := vv.(map[interface{}]interface{})
			method["tags"] = []string{path.(string)}
			methods[m] = method
		}
	}

	cacheDir := config.GetCacheDirectory(g.basePath)
	fmt.Printf("debug: cache dir: %s\n", cacheDir)
	targetDir := cacheDir + "/" + g.schemaDir

	err = copy.Copy(filepath.Join(g.basePath, g.schemaDir), targetDir, copy.Options{
		OnDirExists: func(src, dest string) copy.DirExistsAction {
			return copy.Replace
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to prepare temp api schema: %w", err)
	}

	targetPath := targetDir + "/api.yaml"
	f, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to create api yaml: %w", err)
	}

	err = yaml.NewEncoder(f).Encode(doc)
	if err != nil {
		return "", fmt.Errorf("failed to create api yaml: %w", err)
	}

	return targetPath, nil
}

func (g *OpenAPIGenerator) doGenerate(ctx context.Context, schemaPath string, targetPath string) error {
	langStr := string(GENERATOR_LANGUAGE_GO)
	path, err := config.DumpAssets(g.basePath, "openapi/server-template/"+langStr, "openapi/server-template")
	if err != nil {
		return fmt.Errorf("failed to dump assets: %w", err)
	}
	fmt.Printf("debug: dumped path: %s\n", path)

	generatedPath := filepath.Join(g.basePath, targetPath, "generated")

	err = runOpenapiGenerator(ctx, g.basePath, schemaPath, path, generatedPath, g.info)
	if err != nil {
		return fmt.Errorf("failed to run openapi-generator: %w", err)
	}

	apiPath := filepath.Join(generatedPath, "api")
	err = sanitizeHandlersImports(apiPath)
	if err != nil {
		return err
	}

	handlersPath := filepath.Join(g.basePath, targetPath, "handlers")
	err = moveHandlers(apiPath, handlersPath)
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
func formatGenerated(apiPath string) error {
	return filepath.WalkDir(apiPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		f, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		fmtData, err := format.Source(data)
		if err != nil {
			return err
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
		_, err = w.Write(fmtData)
		if err != nil {
			return err
		}

		err = w.Flush()
		if err != nil {
			return err
		}

		return nil
	})
}

// FIXME: go-specific
func moveHandlers(apiPath string, handlersPath string) error {
	services, err := filepath.Glob(filepath.Join(apiPath, "api_*_service.go"))
	if err != nil {
		return err
	}
	fmt.Printf("services: %v\n", services)
	if len(services) == 0 {
		fmt.Printf("debug: no handlers to move\n")
		return nil
	}
	for _, service := range services {
		baseName := filepath.Base(service)
		baseName = strings.TrimPrefix(baseName, "api_")
		baseName = strings.TrimSuffix(baseName, "_service.go")
		baseName = strings.ReplaceAll(baseName, "_", "/")

		fmt.Printf("debug: processing handler for: %v\n", baseName)
		targetFile := filepath.Join(handlersPath, baseName, "service.go")
		defer func(svc string) {
			if err := os.Remove(svc); err != nil {
				fmt.Printf("failed to remove service file: %s: %s\n", svc, err)
				return
			}
			fmt.Printf("debug: cleaned generated service file: %s\n", svc)
		}(service)

		if _, err := os.Stat(targetFile); err == nil {
			fmt.Printf("debug: skipping existing handler for: %v\n", baseName)
			continue
		}
		if err := os.MkdirAll(filepath.Join(handlersPath, baseName), 0755); err != nil {
			return err
		}
		if err := createHandlersFile(service, targetFile); err != nil {
			return err
		}
		fmt.Printf("debug: created handler for: %v\n", baseName)
	}
	return nil
}

// FIXME: go-specific
func createHandlersFile(serviceFile string, targetFile string) error {
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
		fmt.Printf("debug: writing param %s %s\n", decl.Name, decl.Type)
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
func sanitizeHandlersImports(apiPath string) error {
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

	fmt.Printf("debug: sanitized routes imports\n")
	return nil
}

func runOpenapiGenerator(ctx context.Context, basePath string, schemaPath string, templatePath string, targetDir string,
	info OpenAPIGeneratorInfo) error {
	const (
		image = "openapitools/openapi-generator-cli:v5.3.0"
	)
	curUser, err := user.Current()
	if err != nil {
		return err
	}

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}

	err = copyFile(
		filepath.Join(templatePath, "ignore-list.txt"),
		filepath.Join(targetDir, ".openapi-generator-ignore"),
	)
	if err != nil {
		return err
	}

	templatePathRel := strings.Replace(templatePath, basePath, "", 1)
	schemaPathRel := strings.Replace(schemaPath, basePath, "", 1)
	targetDirRel := strings.Replace(targetDir, basePath, "", 1)
	args := []string{
		"generate",
		"-c", filepath.Join("/repo", templatePathRel, "config.yaml"),
		"-i", filepath.Join("/repo", schemaPathRel),
		"-o", filepath.Join("/repo", targetDirRel),
		"--git-host", info.GitHost,
		"--git-user-id", info.GitNamespace,
		"--git-repo-id", info.GitRepository,
		"--group-id", info.ServiceName,
		"--artifact-id", info.ServiceName,
	}
	fmt.Printf("debug: running docker %s\n", args)

	fmt.Printf("running openapi-generator\n")
	params := docker.DockerRunParams{
		User:   curUser,
		Mounts: map[string]string{"/repo": basePath},
		Cmd:    args,
	}
	err = docker.Run(ctx, image, params)
	if err != nil {
		return err
	}
	fmt.Printf("debug: generated openapi\n")

	return nil
}

func copyFile(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0644)
}
