package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/service/generate"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

const (
	apiSchemaPath = "schemas/%s/api"
	apiServicePath = "go_services/internal/%s"
	svcLanguage generate.GeneratorLanguage = generate.GENERATOR_LANGUAGE_GO
)

func Generate(workspaceContext workspace.Context, name string) error {
	// check if service exists

	if err := generateServiceOpenAPI(workspaceContext.Config, workspaceContext.BasePath, name); err != nil {
		return err
	}
	return nil
}

func generateServiceOpenAPI(conf config.WorkspaceConfig, basePath string, name string) error {
	schemaPath := fmt.Sprintf(apiSchemaPath, name)
	if _, err := os.Stat(filepath.Join(basePath, schemaPath)); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("debug: skipping openapi generating, schema not found for service: %s\n", name)
		return nil
	}
	info := generate.OpenAPIGeneratorInfo{
		GitHost: conf.GitHost,
		GitNamespace: conf.GitNamespace,
		GitRepository: filepath.Join(conf.GitRepository, "go_services"),
		ServiceName: name,
	}
	openapigen := generate.NewOpenAPIGenerator(basePath, schemaPath, svcLanguage, info)
	return openapigen.GenerateServer(fmt.Sprintf(apiServicePath, name))
}
