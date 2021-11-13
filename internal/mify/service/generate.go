package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

var (
	ErrSkip = errors.New("skip step")
)

func Generate(ctx *core.Context, workspaceContext workspace.Context, name string) error {
	serviceConf, err := config.ReadServiceConfig(workspaceContext.BasePath, name)
	if err != nil {
		return err
	}
	repo := fmt.Sprintf("%s/%s/%s",
		workspaceContext.Config.GitHost,
		workspaceContext.Config.GitNamespace,
		workspaceContext.Config.GitRepository)
	tcontext := Context{
		ServiceName: name,
		Repository:  repo,
		GoModule:    repo + "/go_services",
		Workspace:   workspaceContext,
	}

	if err := generateServiceOpenAPI(ctx, tcontext, serviceConf, name); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
	return nil
}

func generateServiceOpenAPI(ctx *core.Context, serviceCtx Context, serviceConf config.ServiceConfig, name string) error {
	hasGenerateTasks := false
	defer func() {
		if !hasGenerateTasks {
			fmt.Println("Nothing to do")
		}
	}()

	err := checkOpenAPISchemas(serviceCtx.Workspace.BasePath, serviceConf, name)
	if err != nil {
		if err == ErrSkip {
			return nil
		}
		return err
	}

	pool, err := util.NewJobPool(ctx, config.GetCacheDirectory(serviceCtx.Workspace.BasePath), runtime.NumCPU())
	if err != nil {
		return err
	}
	defer pool.ClosePool()

	clDiff, err := generateClientsContextStep(ctx, pool, serviceCtx, serviceConf)
	if err != nil && err != ErrSkip {
		return err
	}
	if err == nil {
		hasGenerateTasks = true
	}

	err = generateOpenAPIGeneratorStep(ctx, pool, serviceCtx, serviceConf, name, clDiff)
	if err != nil && err != ErrSkip {
		return err
	}
	if err == nil {
		hasGenerateTasks = true
	}

	return nil
}

func checkOpenAPISchemas(basePath string, serviceConf config.ServiceConfig, name string) error {
	schemaPath := fmt.Sprintf(apiSchemaPath, name)
	if _, err := os.Stat(filepath.Join(basePath, schemaPath)); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("debug: skipping openapi generating, schema not found for service: %s\n", name)
		return ErrSkip
	}

	for clientName := range serviceConf.OpenAPI.Clients {
		clientSchemaPath := fmt.Sprintf(apiSchemaPath, clientName)
		if _, err := os.Stat(filepath.Join(basePath, clientSchemaPath)); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("schema not found while generating client for: %s", clientName)
		}
	}

	return nil
}

func generateOpenAPIGeneratorStep(ctx *core.Context, pool *util.JobPool, serviceCtx Context, serviceConf config.ServiceConfig, name string, clientsDiff clientsDiff) error {
	schemaPath := fmt.Sprintf(apiSchemaPath, name)

	info := generate.OpenAPIGeneratorInfo{
		GitHost: serviceCtx.Workspace.Config.GitHost,
		GitNamespace: serviceCtx.Workspace.Config.GitNamespace,
		GitRepository: serviceCtx.Repository,
		GoModule: serviceCtx.GoModule,
		ServiceName: name,
	}
	targetDir := fmt.Sprintf(apiServicePath, name)

	openapigen := generate.NewOpenAPIGenerator(serviceCtx.Workspace.BasePath, svcLanguage, info)

	isAnyClientGenerating := false
	clientCheckMap := map[string]bool{}
	for clientName := range serviceConf.OpenAPI.Clients {
		clientSchemaPath := fmt.Sprintf(apiSchemaPath, clientName)

		needGenerateClient, err := openapigen.NeedGenerateClient(ctx, clientSchemaPath)
		if err != nil {
			return err
		}

		_, isAdded := clientsDiff.added[clientName]
		if isAdded {
			needGenerateClient = true
		}

		clientCheckMap[clientName] = needGenerateClient
		if needGenerateClient {
			isAnyClientGenerating = true
		}
	}

	needGenerateServer, err := openapigen.NeedGenerateServer(ctx, schemaPath)
	if err != nil {
		return err
	}

	for clientName, _ := range clientsDiff.removed {
		err := openapigen.RemoveClient(ctx, clientName, targetDir)
		if err != nil {
			return err
		}
	}

	if !needGenerateServer && !isAnyClientGenerating {
		return ErrSkip
	}

	err = openapigen.Prepare(pool)
	if err != nil {
		return err
	}

	if needGenerateServer {
		pool.AddJob(util.Job{
			Name: "generate:server",
			Func: func(ctx *core.Context) error {
				if err := openapigen.GenerateServer(ctx, schemaPath, targetDir); err != nil {
					return err
				}
				return nil
			},
		})
	}


	for clientName := range serviceConf.OpenAPI.Clients {
		if !clientCheckMap[clientName] {
			continue
		}
		clientSchemaPath := fmt.Sprintf(apiSchemaPath, clientName)

		client := clientName
		pool.AddJob(util.Job{
			Name: "generate:"+client,
			Func: func(ctx *core.Context) error {
				if err := openapigen.GenerateClient(ctx, client, clientSchemaPath, targetDir); err != nil {
					return fmt.Errorf("failed to generate client for: %s: %w", client, err)
				}
				return nil
			},
		})
	}

	jerr := pool.Run()
	if jerr != nil {
		util.ShowJobError(pool, jerr)
		return jerr.Err
	}

	return nil
}
