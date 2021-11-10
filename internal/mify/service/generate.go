package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/service/generate"
	"github.com/chebykinn/mify/internal/mify/util"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

const (
	apiSchemaPath = "schemas/%s/api"
	apiServicePath = "go_services/internal/%s"
	svcLanguage generate.GeneratorLanguage = generate.GENERATOR_LANGUAGE_GO
)

func Generate(ctx *core.Context, workspaceContext workspace.Context, name string) error {
	serviceConf, err := config.ReadServiceConfig(workspaceContext.BasePath, name)
	if err != nil {
		return err
	}

	if err := generateServiceOpenAPI(ctx, workspaceContext.Config, serviceConf, workspaceContext.BasePath, name); err != nil {
		return err
	}
	return nil
}

func generateServiceOpenAPI(ctx *core.Context, conf config.WorkspaceConfig, serviceConf config.ServiceConfig, basePath string, name string) error {
	schemaPath := fmt.Sprintf(apiSchemaPath, name)
	if _, err := os.Stat(filepath.Join(basePath, schemaPath)); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("debug: skipping openapi generating, schema not found for service: %s\n", name)
		return nil
	}
	fmt.Printf("running generate for %s\n", name)

	info := generate.OpenAPIGeneratorInfo{
		GitHost: conf.GitHost,
		GitNamespace: conf.GitNamespace,
		GitRepository: filepath.Join(conf.GitRepository, "go_services"),
		ServiceName: name,
	}

	pool, err := util.NewJobPool(ctx, config.GetCacheDirectory(basePath), 4)
	if err != nil {
		return err
	}

	openapigen, err := generate.NewOpenAPIGenerator(pool, basePath, svcLanguage, info)
	if err != nil {
		return err
	}

	pool.AddJob(util.Job{
		Name: "generate:server",
		Func: func(ctx *core.Context) error {
			if err := openapigen.GenerateServer(ctx, schemaPath, fmt.Sprintf(apiServicePath, name)); err != nil {
				return err
			}

			return nil
		},
	})


	for clientName := range serviceConf.OpenAPI.Clients {
		clientSchemaPath := fmt.Sprintf(apiSchemaPath, clientName)
		if _, err := os.Stat(filepath.Join(basePath, clientSchemaPath)); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("schema not found while generating client for: %s", clientName)
		}

		client := clientName
		pool.AddJob(util.Job{
			Name: "generate:"+client,
			Func: func(ctx *core.Context) error {
				if err := openapigen.GenerateClient(ctx, client, clientSchemaPath, fmt.Sprintf(apiServicePath, name)); err != nil {
					return fmt.Errorf("failed to generate client for: %s: %w", client, err)
				}
				return nil
			},
		})
	}

	jerr := pool.Run()
	defer pool.ClosePool()

	if jerr != nil && errors.Is(jerr.Err, context.Canceled) {
		return nil
	}

	if jerr != nil {
		fmt.Printf("task %s error: %s\n", jerr.Name, jerr.Err)
		logFile, err := os.Open(pool.GetJobLogPath(jerr.Name))
		if err != nil {
			fmt.Printf("failed to read job %s log: %s", jerr.Name, err)
		}
		fmt.Printf("\nfull log:\n")
		io.Copy(os.Stderr, logFile)
		logFile.Close()

		return jerr.Err
	}

	return nil
}
