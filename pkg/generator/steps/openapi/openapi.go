package openapi

import (
	"bufio"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"

	"os/user"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mify-io/mify/internal/mify/config"
	"github.com/mify-io/mify/internal/mify/util/docker"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/util/threading"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

const (
	CACHE_SERVER_SUBDIR = "server"
	CACHE_CLIENT_SUBDIR = "client"
	SERVER_PACKAGE_NAME = "openapi"

	FILE_TIME_FILENAME     = ".timestamps.yaml"
	GENERATED_API_FILENAME = "api_generated.yaml"
)

type fileTimeMap map[string]int64

type OpenAPIGeneratorInfo struct {
	GitHost       string
	GitNamespace  string
	GitRepository string
	GoModule      string
	ServiceName   string
}

type OpenAPIGenerator struct {
	basePath         string
	language         mifyconfig.ServiceLanguage
	info             OpenAPIGeneratorInfo
	serverAssetsPath string
	clientAssetsPath string

	prepareMutex sync.Mutex
	prepared     bool
}

type OpenAPIGeneratorMode int16

const (
	Server OpenAPIGeneratorMode = iota
	Client
)

func NewOpenAPIGenerator(ctx *gencontext.GenContext) OpenAPIGenerator {
	info := OpenAPIGeneratorInfo{
		GitHost:       ctx.GetWorkspace().Config.GitHost,
		GitNamespace:  ctx.GetWorkspace().Config.GitNamespace,
		GitRepository: ctx.GetWorkspace().Config.GitRepository,
		GoModule:      ctx.GetWorkspace().GetGoModule(),
		ServiceName:   ctx.GetServiceName(),
	}

	return OpenAPIGenerator{
		basePath: ctx.GetWorkspace().BasePath,
		language: ctx.MustGetMifySchema().Language,
		info:     info,
	}
}

func (g *OpenAPIGenerator) PrepareSync(ctx *gencontext.GenContext) error {
	if g.prepared {
		return nil
	}

	return threading.DoUnderLock(&g.prepareMutex, func() error {
		if g.prepared {
			return nil
		}

		return g.Prepare(ctx)
	})
}

func (g *OpenAPIGenerator) Prepare(ctx *gencontext.GenContext) error {
	const (
		image = "openapitools/openapi-generator-cli:v5.3.0"
	)

	langStr := string(g.language)
	// TODO: pool?
	if err := docker.Cleanup(ctx.GetGoContext()); err != nil {
		return err
	}

	logFile, err := createLogFile(ctx, "docker-pull.log")
	if err != nil {
		return err
	}
	defer logFile.Close()

	if err := docker.PullImage(ctx.GetGoContext(), ctx.Logger, logFile, image); err != nil {
		return err
	}

	if config.HasAssets("openapi/server-template/" + langStr) {
		g.serverAssetsPath, err = config.DumpAssets(
			ctx.GetWorkspace().GetCacheDirectory(),
			"openapi/server-template/"+langStr, "openapi/server-template")
		if err != nil {
			return fmt.Errorf("failed to dump assets: %w", err)
		}
		ctx.Logger.Infof("dumped server path: %s", g.serverAssetsPath)
	}

	if config.HasAssets("openapi/client-template/" + langStr) {
		g.clientAssetsPath, err = config.DumpAssets(
			ctx.GetWorkspace().GetCacheDirectory(),
			"openapi/client-template/"+langStr, "openapi/client-template")
		if err != nil {
			return fmt.Errorf("failed to dump assets: %w", err)
		}
		ctx.Logger.Infof("dumped client path: %s", g.clientAssetsPath)
	}

	g.prepared = true
	return nil
}

func (g *OpenAPIGenerator) NeedGenerateServer(ctx *gencontext.GenContext, schemaRelDir string) (bool, error) {
	return g.anySchemaChanged(ctx, Server, schemaRelDir)
}

func (g *OpenAPIGenerator) NeedGenerateClient(ctx *gencontext.GenContext, schemaRelDir string) (bool, error) {
	return g.anySchemaChanged(ctx, Client, schemaRelDir)
}

func (g *OpenAPIGenerator) GenerateServer(ctx *gencontext.GenContext, outputDir string) error {
	schemaDir := ctx.GetWorkspace().GetApiSchemaDirRelPath(ctx.GetServiceName())
	if len(g.serverAssetsPath) == 0 {
		return fmt.Errorf("failed to generate server: no generator available for language: %s", g.language)
	}
	schemaPath, paths, err := g.makeServerEnrichedSchema(ctx, schemaDir)
	if err != nil {
		return fmt.Errorf("failed to generate server: %w", err)
	}

	err = g.doGenerateServer(ctx, g.serverAssetsPath, schemaPath, outputDir, paths)
	if err != nil {
		return fmt.Errorf("failed to generate server: %w", err)
	}

	err = updateGenerationTime(ctx, ctx.GetServiceName(), filepath.Dir(schemaPath))
	if err != nil {
		return err
	}

	return nil

}

func (g *OpenAPIGenerator) GenerateClient(ctx *gencontext.GenContext, clientName string, outputDir string) error {
	if len(g.clientAssetsPath) == 0 {
		return fmt.Errorf("failed to generate client: no generator available for language: %s", g.language)
	}
	inputSchemaPath := ctx.GetWorkspace().GetApiSchemaAbsPath(clientName, "api.yaml")
	schemaPath, err := g.makeClientEnrichedSchema(ctx, inputSchemaPath)
	if err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	err = g.doGenerateClient(ctx, g.clientAssetsPath, clientName, schemaPath, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	err = updateGenerationTime(ctx, clientName, filepath.Dir(schemaPath))
	if err != nil {
		return err
	}

	return nil

}

func (g *OpenAPIGenerator) RemoveClient(ctx *gencontext.GenContext, clientName string, outputDir string) error {
	return g.doRemoveClient(ctx, clientName, outputDir)
}

// private

func (g *OpenAPIGenerator) readSchema(ctx *gencontext.GenContext, schemaPath string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}
	loader := &openapi3.Loader{
		Context:               ctx.GetGoContext(),
		IsExternalRefsAllowed: true,
	}
	url, err := url.Parse(schemaPath)
	if err != nil {
		return nil, err
	}

	openapiDoc, err := loader.LoadFromDataWithPath(data, url)
	if err != nil {
		return nil, err
	}
	err = openapiDoc.Validate(ctx.GetGoContext())
	if err != nil {
		return nil, err
	}

	doc := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (g *OpenAPIGenerator) getTempSchemaPath(cacheDir string, schemaPath string, cacheSubdir string) string {
	schemaDir := filepath.Dir(strings.Replace(schemaPath, g.basePath, "", 1))
	return filepath.Join(cacheDir, schemaDir, cacheSubdir)
}

func (g *OpenAPIGenerator) saveEnrichedSchema(
	ctx *gencontext.GenContext, doc map[string]interface{}, schemaPath string, cacheSubdir string) (string, error) {
	schemaDir := filepath.Dir(strings.Replace(schemaPath, g.basePath, "", 1))
	cacheDir := ctx.GetWorkspace().GetCacheDirectory()
	targetDir := filepath.Join(cacheDir, schemaDir, cacheSubdir)
	ctx.Logger.Infof("saving schema in: %s", targetDir)

	err := copy.Copy(filepath.Join(g.basePath, schemaDir), targetDir, copy.Options{
		OnDirExists: func(src, dest string) copy.DirExistsAction {
			return copy.Replace
		},
		Skip: func(src string) (bool, error) {
			return filepath.Base(src) == GENERATED_API_FILENAME, nil
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

// FIXME: go-specific
func formatGenerated(apiPath string, language mifyconfig.ServiceLanguage) error {
	if language != mifyconfig.ServiceLanguageGo {
		return nil
	}
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

		data, err := ioutil.ReadAll(f)
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

func makeFileUpdateMap(ctx *gencontext.GenContext, schemaDir string, tmpSchemaDir string) (fileTimeMap, error) {
	basePath := ctx.GetWorkspace().BasePath

	fileMap := fileTimeMap{}
	err := filepath.WalkDir(filepath.Join(basePath, schemaDir), func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".yaml" {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		fileMap[path] = info.ModTime().UnixNano()
		return nil
	})
	if err != nil {
		return nil, err
	}

	return fileMap, nil
}

func updateGenerationTime(ctx *gencontext.GenContext, targetServiceName string, tmpSchemaDir string) error {
	schemaDir := ctx.GetWorkspace().GetApiSchemaDirRelPath(targetServiceName)
	ctx.Logger.Infof("updating generation time in: %s", schemaDir)

	fileMap, err := makeFileUpdateMap(ctx, schemaDir, tmpSchemaDir)
	if err != nil {
		return fmt.Errorf("failed to write file update times: %w", err)
	}

	f, err := os.OpenFile(filepath.Join(tmpSchemaDir, FILE_TIME_FILENAME), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file update times: %w", err)
	}

	err = yaml.NewEncoder(f).Encode(fileMap)
	if err != nil {
		return fmt.Errorf("failed to write file update times: %w", err)
	}
	return nil
}

func (g *OpenAPIGenerator) anySchemaChanged(ctx *gencontext.GenContext, mode OpenAPIGeneratorMode, schemaRelDir string) (bool, error) {
	fullPath := filepath.Join(g.basePath, schemaRelDir)
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return false, err
	}

	for _, f := range files {
		cache_subdir := CACHE_SERVER_SUBDIR
		if mode == Client {
			cache_subdir = CACHE_CLIENT_SUBDIR
		}

		schemaFile := filepath.Join(fullPath, f.Name())
		schemaPath := g.getTempSchemaPath(ctx.GetWorkspace().GetCacheDirectory(), schemaFile, cache_subdir)

		changed, err := isSchemasChanged(ctx, g.basePath, schemaRelDir, schemaPath)
		if err != nil {
			return false, err
		}

		if changed {
			return true, nil
		}
	}

	return false, nil
}

func isSchemasChanged(ctx *gencontext.GenContext, basePath string, schemaDir string, tmpSchemaDir string) (bool, error) {
	f, err := os.Open(filepath.Join(tmpSchemaDir, FILE_TIME_FILENAME))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}
		return false, fmt.Errorf("failed to compare file update times: %w", err)
	}

	var oldMap fileTimeMap
	err = yaml.NewDecoder(f).Decode(&oldMap)
	if err != nil {
		return false, fmt.Errorf("failed to compare file update times: %w", err)
	}

	fileMap, err := makeFileUpdateMap(ctx, schemaDir, tmpSchemaDir)
	if err != nil {
		return false, fmt.Errorf("failed to compare file update times: %w", err)
	}

	if len(fileMap) != len(oldMap) {
		return true, nil
	}

	for path, time := range oldMap {
		newTime, ok := fileMap[path]
		if !ok {
			return true, nil
		}
		if newTime != time {
			return true, nil
		}
	}

	return false, nil
}

// TODO: refactor mixed client/service generation
func runOpenapiGenerator(
	ctx *gencontext.GenContext, basePath string, schemaPath string, templatePath string, targetDir string,
	packageName string,
	clientName string,
	servicePort int,
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
		"-p", "goModule=" + info.GoModule,
		"-p", "serviceName=" + info.ServiceName,
		"-p", "clientName=" + clientName,
		"-p", "clientEndpointEnv=" + MakeClientEnvName(clientName),
		"-p", "serviceEndpointEnv=" + MakeServerEnvName(info.ServiceName),
		"-p", "serviceEndpoint=:" + strconv.Itoa(servicePort),
		"--package-name", packageName,
	}
	ctx.Logger.Infof("running docker %s", args)

	ctx.Logger.Infof("running openapi-generator")
	params := docker.DockerRunParams{
		User:   curUser,
		Mounts: map[string]string{"/repo": basePath},
		Cmd:    args,
	}

	logFileName := fmt.Sprintf("openapi-generator-run-%s.log", clientName)
	logFile, err := createLogFile(ctx, logFileName)
	if err != nil {
		return err
	}
	// TODO: move globally
	defer func() {
		logFile.Close()
		if err == nil {
			return
		}
		ctx.Logger.Errorf("openapi-generator task failed, dumping last errors, see full logs in: %s", logFile.Name())
		file, err := os.Open(logFile.Name())
		if err != nil {
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return
		}
	}()

	err = docker.Run(ctx.GetGoContext(), ctx.Logger, logFile, image, params)
	if err != nil {
		return err
	}
	ctx.Logger.Infof("generated openapi")

	return nil
}

func copyFile(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0644)
}

func createLogFile(ctx *gencontext.GenContext, fileName string) (*os.File, error) {
	// TODO: library for creating log files
	logsDir := ctx.GetWorkspace().GetLogsDirectory()
	logFile := path.Join(logsDir, fmt.Sprintf("%s-%s", ctx.GetServiceName(), fileName))

	f, err := os.OpenFile(logFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	return f, nil
}
