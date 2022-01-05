package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/service/generate"
	"github.com/chebykinn/mify/internal/mify/service/lang"
	"github.com/chebykinn/mify/internal/mify/util"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

const (
	apiSchemaPath                      = "schemas/%s/api"
	svcLanguage   lang.ServiceLanguage = lang.ServiceLanguageGo
)

var (
	ErrSkip = errors.New("skip step")
)

func getAPIServicePathByLang(language lang.ServiceLanguage, serviceName string) (string, error) {
	switch language {
	case lang.ServiceLanguageGo:
		return "go_services/internal/" + serviceName, nil
	case lang.ServiceLanguageJs:
		return "js_services/" + serviceName, nil
	}
	return "", fmt.Errorf("unknown language: %s", language)
}

func Generate(ctx *core.Context, pool *util.JobPool, workspaceContext workspace.Context, name string) error {
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
		Language:    serviceConf.Language,
		GoModule:    repo + "/go_services",
		Workspace:   workspaceContext,
	}

	if err := generateServiceOpenAPI(ctx, pool, tcontext, serviceConf, name); err != nil {
		return err
	}

	// TODO: suppport dev-runner for all languages?
	if serviceConf.Language == lang.ServiceLanguageGo {
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

func generateServiceOpenAPI(ctx *core.Context, pool *util.JobPool, serviceCtx Context, serviceConf config.ServiceConfig, name string) error {
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

func checkOpenAPISchemas(ctx *core.Context, basePath string, serviceConf config.ServiceConfig, name string) (bool, error) {
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

func generateOpenAPIGeneratorStep(ctx *core.Context, pool *util.JobPool, serviceCtx Context, serviceConf config.ServiceConfig, hasServer bool, name string, clientsDiff clientsDiff) error {
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
