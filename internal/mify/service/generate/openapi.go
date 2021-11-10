package generate

import (
	"bufio"
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
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/util"
	"github.com/chebykinn/mify/internal/mify/util/docker"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

type GeneratorLanguage string

const (
	GENERATOR_LANGUAGE_GO GeneratorLanguage = "go"

	CACHE_SERVER_SUBDIR = "server"
	CACHE_CLIENT_SUBDIR = "client"
	SERVER_PACKAGE_NAME = "openapi"
)

type OpenAPIGeneratorInfo struct {
	GitHost       string
	GitNamespace  string
	GitRepository string
	ServiceName   string
}

type OpenAPIGenerator struct {
	basePath  string
	language  GeneratorLanguage
	info      OpenAPIGeneratorInfo
}

func NewOpenAPIGenerator(pool *util.JobPool, basePath string, language GeneratorLanguage, info OpenAPIGeneratorInfo) (OpenAPIGenerator, error) {
	const (
		image = "openapitools/openapi-generator-cli:v5.3.0"
	)

	pool.AddJob(util.Job{
		Name:"generate:prepare",
		Func: func(ctx *core.Context) error {
			if err := docker.Cleanup(ctx.Ctx); err != nil {
				return err
			}
			if err := docker.PullImage(ctx.Ctx, ctx.Logger, image); err != nil {
				return err
			}
			return nil
		},
	})

	jerr := pool.Run()
	if jerr != nil {
		return OpenAPIGenerator{}, jerr.Err
	}

	return OpenAPIGenerator{
		basePath:  basePath,
		language:  language,
		info:      info,
	}, nil
}

func (g *OpenAPIGenerator) GenerateServer(ctx *core.Context, schemaDir string, outputDir string) error {
	inputSchemaPath := filepath.Join(g.basePath, schemaDir, "/api.yaml")
	schemaPath, err := g.makeServerEnrichedSchema(ctx, inputSchemaPath)
	if err != nil {
		return fmt.Errorf("failed to generate server: %w", err)
	}

	err = g.doGenerateServer(ctx, schemaPath, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate server: %w", err)
	}

	return nil

}

func (g *OpenAPIGenerator) GenerateClient(ctx *core.Context, clientName string, schemaDir string, outputDir string) error {
	inputSchemaPath := filepath.Join(g.basePath, schemaDir, "/api.yaml")
	schemaPath, err := g.makeClientEnrichedSchema(ctx, inputSchemaPath)
	if err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	err = g.doGenerateClient(ctx, clientName, schemaPath, outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	return nil

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
func formatGenerated(apiPath string) error {
	return filepath.WalkDir(apiPath, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return fmt.Errorf("failed to format: %s: %w", apiPath, err)
		}
		if d.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != "go" {
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

func runOpenapiGenerator(
	ctx *core.Context, basePath string, schemaPath string, templatePath string, targetDir string,
	packageName string,
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
		"--package-name", packageName,
		"--git-host", info.GitHost,
		"--git-user-id", info.GitNamespace,
		"--git-repo-id", info.GitRepository,
		"--group-id", info.ServiceName,
		"--artifact-id", info.ServiceName,
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
