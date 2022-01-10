package generate

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
	"strconv"

	"os/user"
	"path/filepath"
	"strings"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/util"
	"github.com/chebykinn/mify/internal/mify/util/docker"
	"github.com/chebykinn/mify/pkg/mifyconfig"
	"github.com/getkin/kin-openapi/openapi3"
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
}

type OpenAPIGeneratorMode int16

const (
	Server OpenAPIGeneratorMode = iota
	Client
)

func NewOpenAPIGenerator(basePath string, language mifyconfig.ServiceLanguage, info OpenAPIGeneratorInfo) OpenAPIGenerator {
	return OpenAPIGenerator{
		basePath: basePath,
		language: language,
		info:     info,
	}
}

func (g *OpenAPIGenerator) Prepare(pool *util.JobPool) error {
	const (
		image = "openapitools/openapi-generator-cli:v5.3.0"
	)

	langStr := string(g.language)
	pool.AddJob(util.Job{
		Name: "generate:prepare",
		Func: func(ctx *core.Context) error {
			if err := docker.Cleanup(ctx.Ctx); err != nil {
				return err
			}
			if err := docker.PullImage(ctx.Ctx, ctx.Logger, image); err != nil {
				return err
			}
			var err error
			if config.HasAssets("openapi/server-template/" + langStr) {
				g.serverAssetsPath, err = config.DumpAssets(
					g.basePath, "openapi/server-template/"+langStr, "openapi/server-template")
				if err != nil {
					return fmt.Errorf("failed to dump assets: %w", err)
				}
				ctx.Logger.Printf("dumped server path: %s\n", g.serverAssetsPath)
			}

			if config.HasAssets("openapi/client-template/" + langStr) {
				g.clientAssetsPath, err = config.DumpAssets(
					g.basePath, "openapi/client-template/"+langStr, "openapi/client-template")
				if err != nil {
					return fmt.Errorf("failed to dump assets: %w", err)
				}
				ctx.Logger.Printf("dumped client path: %s\n", g.clientAssetsPath)
			}
			return nil
		},
	})

	if err := pool.Run(); err != nil {
		return err
	}

	return nil
}

func (g *OpenAPIGenerator) NeedGenerateServer(ctx *core.Context, schemaDir string) (bool, error) {
	return g.anySchemaChanged(ctx, Server, schemaDir)
}

func (g *OpenAPIGenerator) NeedGenerateClient(ctx *core.Context, schemaDir string) (bool, error) {
	return g.anySchemaChanged(ctx, Client, schemaDir)
}

func (g *OpenAPIGenerator) GenerateServer(ctx *core.Context, schemaDir string, outputDir string) error {
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

	err = updateGenerationTime(ctx, g.basePath, schemaDir, filepath.Dir(schemaPath))
	if err != nil {
		return err
	}

	return nil

}

func (g *OpenAPIGenerator) GenerateClient(ctx *core.Context, clientName string, schemaDir string, outputDir string) error {
	if len(g.clientAssetsPath) == 0 {
		return fmt.Errorf("failed to generate client: no generator available for language: %s", g.language)
	}
	inputSchemaPath := filepath.Join(g.basePath, schemaDir, "/api.yaml")
	schemaPath, err := g.makeClientEnrichedSchema(ctx, inputSchemaPath)
	if err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	err = g.doGenerateClient(ctx, g.clientAssetsPath, clientName, schemaPath, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	err = updateGenerationTime(ctx, g.basePath, schemaDir, filepath.Dir(schemaPath))
	if err != nil {
		return err
	}

	return nil

}

func (g *OpenAPIGenerator) RemoveClient(ctx *core.Context, clientName string, outputDir string) error {
	return g.doRemoveClient(ctx, clientName, outputDir)
}

// private

func (g *OpenAPIGenerator) readSchema(ctx *core.Context, schemaPath string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}
	loader := &openapi3.Loader{
		Context:               ctx.Ctx,
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
	err = openapiDoc.Validate(ctx.Ctx)
	if err != nil {
		return nil, err
	}

	doc := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (g *OpenAPIGenerator) getTempSchemaPath(schemaPath string, cacheSubdir string) string {
	schemaDir := filepath.Dir(strings.Replace(schemaPath, g.basePath, "", 1))
	cacheDir := config.GetCacheDirectory(g.basePath)
	return filepath.Join(cacheDir, schemaDir, cacheSubdir)
}

func (g *OpenAPIGenerator) saveEnrichedSchema(
	ctx *core.Context, doc map[string]interface{}, schemaPath string, cacheSubdir string) (string, error) {
	schemaDir := filepath.Dir(strings.Replace(schemaPath, g.basePath, "", 1))
	cacheDir := config.GetCacheDirectory(g.basePath)
	targetDir := filepath.Join(cacheDir, schemaDir, cacheSubdir)
	ctx.Logger.Printf("saving schema in: %s\n", targetDir)

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

func makeFileUpdateMap(ctx *core.Context, basePath string, schemaDir string, tmpSchemaDir string) (fileTimeMap, error) {
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

func updateGenerationTime(ctx *core.Context, basePath string, schemaDir string, tmpSchemaDir string) error {
	ctx.Logger.Printf("updating generation time in: %s", schemaDir)

	fileMap, err := makeFileUpdateMap(ctx, basePath, schemaDir, tmpSchemaDir)
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

func (g *OpenAPIGenerator) anySchemaChanged(ctx *core.Context, mode OpenAPIGeneratorMode, schemaDir string) (bool, error) {
	fullPath := filepath.Join(g.basePath, schemaDir)
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
		schemaPath := g.getTempSchemaPath(schemaFile, cache_subdir)

		changed, err := isSchemasChanged(ctx, g.basePath, schemaDir, schemaPath)
		if err != nil {
			return false, err
		}

		if changed {
			return true, nil
		}
	}

	return false, nil
}

func isSchemasChanged(ctx *core.Context, basePath string, schemaDir string, tmpSchemaDir string) (bool, error) {
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

	fileMap, err := makeFileUpdateMap(ctx, basePath, schemaDir, tmpSchemaDir)
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

func runOpenapiGenerator(
	ctx *core.Context, basePath string, schemaPath string, templatePath string, targetDir string,
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
		"-p", "clientEndpointEnv=" + MakeClientEnvName(packageName),
		"-p", "serviceEndpointEnv=" + MakeServerEnvName(info.ServiceName),
		"-p", "serviceEndpoint=:" + strconv.Itoa(servicePort),
		"--package-name", packageName,
	}
	ctx.Logger.Printf("running docker %s\n", args)

	ctx.Logger.Printf("running openapi-generator\n")
	params := docker.DockerRunParams{
		User:   curUser,
		Mounts: map[string]string{"/repo": basePath},
		Cmd:    args,
	}
	err = docker.Run(ctx.Ctx, ctx.Logger, image, params)
	if err != nil {
		return err
	}
	ctx.Logger.Printf("generated openapi\n")

	return nil
}

func copyFile(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0644)
}
