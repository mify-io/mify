package service

import (
	"context"
	"errors"
	"fmt"
	"io"
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

func Generate(ctx *core.Context, workspaceContext workspace.Context, name string) error {
	serviceConf, err := config.ReadServiceConfig(workspaceContext.BasePath, name)
	if err != nil {
		return err
	}
	repo := fmt.Sprintf("%s/%s/%s",
		workspaceContext.Config.GitHost,
		workspaceContext.Config.GitNamespace,
		workspaceContext.Config.GitRepository)
	context := Context{
		ServiceName: name,
		Repository:  repo,
		GoModule:    repo + "/go_services",
		Workspace:   workspaceContext,
	}

	if err := generateServiceOpenAPI(ctx, context, serviceConf, name); err != nil {
		return err
	}
	return nil
}

func generateServiceOpenAPI(ctx *core.Context, serviceCtx Context, serviceConf config.ServiceConfig, name string) error {
	schemaPath := fmt.Sprintf(apiSchemaPath, name)
	if _, err := os.Stat(filepath.Join(serviceCtx.Workspace.BasePath, schemaPath)); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("debug: skipping openapi generating, schema not found for service: %s\n", name)
		return nil
	}
	fmt.Printf("Running code generation for %s\n", name)

	info := generate.OpenAPIGeneratorInfo{
		GitHost: serviceCtx.Workspace.Config.GitHost,
		GitNamespace: serviceCtx.Workspace.Config.GitNamespace,
		GitRepository: serviceCtx.Repository,
		GoModule: serviceCtx.GoModule,
		ServiceName: name,
	}

	pool, err := util.NewJobPool(ctx, config.GetCacheDirectory(serviceCtx.Workspace.BasePath), runtime.NumCPU())
	if err != nil {
		return err
	}


	openapigen := generate.NewOpenAPIGenerator(serviceCtx.Workspace.BasePath, svcLanguage, info)

	isAnyClientGenerating := false
	clientCheckMap := map[string]bool{}
	clientsList := make([]string, 0, len(serviceConf.OpenAPI.Clients))
	clientsCtxList := make([]OpenAPIClientContext, 0, len(serviceConf.OpenAPI.Clients))
	for clientName := range serviceConf.OpenAPI.Clients {
		clientSchemaPath := fmt.Sprintf(apiSchemaPath, clientName)
		if _, err := os.Stat(filepath.Join(serviceCtx.Workspace.BasePath, clientSchemaPath)); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("schema not found while generating client for: %s", clientName)
		}

		needGenerateClient, err := openapigen.NeedGenerateClient(ctx, clientSchemaPath)
		if err != nil {
			return err
		}
		clientCheckMap[clientName] = needGenerateClient
		if needGenerateClient {
			isAnyClientGenerating = true
		}
		packageName := generate.MakePackageName(clientName)
		fieldName := generate.SnakeCaseToCamelCase(generate.SanitizeClientName(clientName), false)
		methodName := generate.SnakeCaseToCamelCase(generate.SanitizeClientName(clientName), true)
		clientsList = append(clientsList, clientName)
		clientsCtxList = append(clientsCtxList, OpenAPIClientContext{
			ClientName: clientName,
			PackageName: packageName,
			PrivateFieldName: fieldName,
			PublicMethodName: methodName,
		})

	}
	serviceCtx.OpenAPI.Clients = clientsCtxList

	needGenerateServer, err := openapigen.NeedGenerateServer(ctx, schemaPath, clientsList)
	if err != nil {
		return err
	}

	hasGenerateTasks := needGenerateServer || isAnyClientGenerating

	if hasGenerateTasks {
		err = openapigen.Prepare(pool)
		if err != nil {
			return err
		}
	}

	if needGenerateServer {
		pool.AddJob(util.Job{
			Name: "generate:server",
			Func: func(ctx *core.Context) error {
				if err := openapigen.GenerateServer(ctx, schemaPath, fmt.Sprintf(apiServicePath, name), clientsList); err != nil {
					return err
				}
				// FIXME: go specific
				// FIXME: regenerate only clients file
				subPath := "go_services/internal/#svc#/generated/core"
				if err := RenderTemplateTreeSubPath(ctx, serviceCtx, subPath); err != nil {
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
				if err := openapigen.GenerateClient(ctx, client, clientSchemaPath, fmt.Sprintf(apiServicePath, name)); err != nil {
					return fmt.Errorf("failed to generate client for: %s: %w", client, err)
				}
				return nil
			},
		})
	}

	var jerr *util.JobError
	if hasGenerateTasks {
		jerr = pool.Run()
		defer pool.ClosePool()
	}

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

	if !hasGenerateTasks {
		fmt.Println("Nothing to do")
	}

	return nil
}
