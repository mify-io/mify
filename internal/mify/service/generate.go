package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/service/apigateway"
	"github.com/chebykinn/mify/internal/mify/service/generate"
	"github.com/chebykinn/mify/internal/mify/util"
	"github.com/chebykinn/mify/internal/mify/workspace"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

const (
	apiSchemaPath                            = "schemas/%s/api"
	svcLanguage   mifyconfig.ServiceLanguage = mifyconfig.ServiceLanguageGo
)

var (
	ErrSkip = errors.New("skip step")
)

func getAPIServicePathByLang(language mifyconfig.ServiceLanguage, serviceName string) (string, error) {
	switch language {
	case mifyconfig.ServiceLanguageGo:
		return mifyconfig.GoServicesRoot + "/internal/" + serviceName, nil
	case mifyconfig.ServiceLanguageJs:
		return mifyconfig.JsServicesRoot + "/" + serviceName, nil
	}
	return "", fmt.Errorf("unknown language: %s", language)
}

func Generate(ctx *core.Context, pool *util.JobPool, workspaceContext workspace.Context, name string) error {
	serviceConf, tcontext, err := initCtx(workspaceContext, name)
	if err != nil {
		return err
	}

	if name == apigateway.ApiGatewayName {
		if err := regenerateApiGateway(ctx, pool, tcontext); err != nil {
			return err
		}

		// Reload configs. They could have changed.
		serviceConf, tcontext, err = initCtx(workspaceContext, name)
		if err != nil {
			return err
		}
	}

	if err := generateServiceOpenAPI(ctx, pool, tcontext, serviceConf, name); err != nil {
		return err
	}

	// TODO: suppport dev-runner for all languages?
	if serviceConf.Language == mifyconfig.ServiceLanguageGo {
		// Hack. We could generate new services during generateServiceOpenAPI, so reload context.
		// TODO: create default service stub in "create service". Or track service list
		// by extra yaml file, without scanning fs.
		tcontext.Workspace, err = workspace.InitContext(workspaceContext.BasePath)
		if err != nil {
			return err
		}

		if err := generateDevRunner(ctx, pool, tcontext); err != nil {
			return err
		}
	}

	return nil
}

func initCtx(workspaceContext workspace.Context, name string) (mifyconfig.ServiceConfig, Context, error) {
	serviceConf, err := mifyconfig.ReadServiceConfig(workspaceContext.BasePath, name)
	if err != nil {
		return mifyconfig.ServiceConfig{}, Context{}, err
	}

	repo := fmt.Sprintf("%s/%s/%s",
		workspaceContext.Config.GitHost,
		workspaceContext.Config.GitNamespace,
		workspaceContext.Config.GitRepository)

	return serviceConf, Context{
		ServiceName: name,
		Repository:  repo,
		Language:    serviceConf.Language,
		GoModule:    repo + "/" + mifyconfig.GoServicesRoot,
		Workspace:   workspaceContext,
	}, nil
}

func generateDevRunner(ctx *core.Context, pool *util.JobPool, serviceCtx Context) error {
	pool.AddJob(util.Job{
		Name: "generate:dev-runner",
		Func: func(ctx *core.Context) error {
			if err := RenderTemplateTreeSubPath(ctx, serviceCtx, "go_services/cmd/dev-runner"); err != nil {
				return err
			}
			return nil
		},
	})
	if err := pool.Run(); err != nil {
		return err
	}
	return nil
}

func regenerateApiGateway(ctx *core.Context, pool *util.JobPool, serviceCtx Context) error {
	var publicApis apigateway.PublicApis = nil
	err := pool.RunImmediate(util.Job{
		Name: "generate:api-gateway-schema",
		Func: func(ctx *core.Context) error {
			var err error
			publicApis, err = apigateway.RegenerateSchema(serviceCtx.Workspace)
			return err
		},
	})

	if err != nil {
		return err
	}

	if publicApis == nil {
		fmt.Println("No public apis were found. Skipping api gateway handlers generation")
		return nil
	}

	err = pool.RunImmediate(util.Job{
		Name: "generate:api-gateway-handlers",
		Func: func(ctx *core.Context) error {
			return apigateway.RegenerateHandlers(ctx, serviceCtx.Workspace, publicApis)
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func generateServiceOpenAPI(ctx *core.Context, pool *util.JobPool, serviceCtx Context, serviceConf mifyconfig.ServiceConfig, name string) error {
	hasGenerateTasks := false
	var err error
	defer func() {
		if !hasGenerateTasks && (err == nil || err == ErrSkip) {
			fmt.Println("Nothing to do")
		}
	}()

	var hasServer bool
	hasServer, err = checkOpenAPISchemas(ctx, serviceCtx.Workspace.BasePath, serviceConf, name)
	if err != nil {
		return err
	}

	var clDiff clientsDiff
	clDiff, err = generateClientsContextStep(ctx, pool, serviceCtx, serviceConf)
	if err != nil && err != ErrSkip {
		return err
	}
	if err == nil {
		hasGenerateTasks = true
	}

	err = generateOpenAPIGeneratorStep(ctx, pool, serviceCtx, serviceConf, hasServer, name, clDiff)
	if err != nil && err != ErrSkip {
		return err
	}
	if err == nil {
		hasGenerateTasks = true
	}

	return nil
}

func checkOpenAPISchemas(ctx *core.Context, basePath string, serviceConf mifyconfig.ServiceConfig, name string) (bool, error) {
	schemaPath := fmt.Sprintf(apiSchemaPath, name)
	hasServer := true
	if _, err := os.Stat(filepath.Join(basePath, schemaPath)); errors.Is(err, os.ErrNotExist) {
		hasServer = false
	}

	for clientName := range serviceConf.OpenAPI.Clients {
		clientSchemaPath := fmt.Sprintf(apiSchemaPath, clientName)
		if _, err := os.Stat(filepath.Join(basePath, clientSchemaPath)); errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("schema not found while generating client for: %s", clientName)
		}
	}

	return hasServer, nil
}

func generateOpenAPIGeneratorStep(ctx *core.Context, pool *util.JobPool, serviceCtx Context, serviceConf mifyconfig.ServiceConfig, hasServer bool, name string, clientsDiff clientsDiff) error {
	schemaPath := fmt.Sprintf(apiSchemaPath, name)

	info := generate.OpenAPIGeneratorInfo{
		GitHost:       serviceCtx.Workspace.Config.GitHost,
		GitNamespace:  serviceCtx.Workspace.Config.GitNamespace,
		GitRepository: serviceCtx.Repository,
		GoModule:      serviceCtx.GoModule,
		ServiceName:   name,
	}
	targetDir, err := getAPIServicePathByLang(serviceCtx.Language, name)
	if err != nil {
		return err
	}

	openapigen := generate.NewOpenAPIGenerator(serviceCtx.Workspace.BasePath, serviceCtx.Language, info)

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

	var needGenerateServer bool
	if hasServer {
		var err error
		needGenerateServer, err = openapigen.NeedGenerateServer(ctx, schemaPath)
		if err != nil {
			return err
		}
	}

	for clientName := range clientsDiff.removed {
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
			Name: "generate:" + client,
			Func: func(ctx *core.Context) error {
				if err := openapigen.GenerateClient(ctx, client, clientSchemaPath, targetDir); err != nil {
					return fmt.Errorf("failed to generate client for: %s: %w", client, err)
				}
				return nil
			},
		})
	}

	if err := pool.Run(); err != nil {
		return err
	}

	return nil
}
