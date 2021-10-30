package generate

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/chebykinn/mify/internal/mify/config"
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

	err = g.doGenerate(schemaPath, outputDir)
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
	// TODO: maybe pass context from caller
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

func (g *OpenAPIGenerator) doGenerate(schemaPath string, targetPath string) error {
	path, err := config.DumpAssets(g.basePath, "openapi/server-template", "openapi")
	if err != nil {
		return err
	}
	fmt.Printf("debug: dumped path: %s\n", path)

	generatedPath := filepath.Join(g.basePath, targetPath, "generated")

	err = runOpenapiGenerator(g.basePath, schemaPath, filepath.Join(path, "server-template"), generatedPath, g.info)
	if err != nil {
		return err
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

	return nil
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
		baseName = strings.ReplaceAll(baseName, "api_", "")
		baseName = strings.ReplaceAll(baseName, "_service.go", "")
		baseName = strings.ReplaceAll(baseName, "_", "/")

		fmt.Printf("debug: processing handler for: %v\n", baseName)
		targetFile := filepath.Join(handlersPath, baseName, "service.go")
		defer func(svc string) {
			if err := os.Remove(svc); err != nil {
				fmt.Printf("failed to remove service file: %s: %s\n", svc, err)
				return
			}
			fmt.Printf("debug: removed service file: %s\n", svc)
		}(service)

		if _, err := os.Stat(targetFile); err == nil {
			fmt.Printf("debug: skipping existing handler for: %v\n", baseName)
			continue
		}
		if err := os.MkdirAll(filepath.Join(handlersPath, baseName), 0755); err != nil {
			return err
		}
		if err := copyFile(service, targetFile); err != nil {
			return err
		}
		fmt.Printf("debug: created handler for: %v\n", baseName)
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

func runOpenapiGenerator(basePath string, schemaPath string, templatePath string, targetDir string,
	info OpenAPIGeneratorInfo) error {
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
		"run",
		"--rm",
		"--user", curUser.Uid + ":" + curUser.Gid,
		"-v", basePath + ":/repo",
		"openapitools/openapi-generator-cli",
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

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	// TODO only if verbose
	fmt.Printf("%s", output)
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
