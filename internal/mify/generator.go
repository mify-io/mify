package mify

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/service"
	"github.com/chebykinn/mify/internal/mify/service/client"
	"github.com/chebykinn/mify/internal/mify/util"
	"github.com/chebykinn/mify/internal/mify/util/docker"
	"github.com/chebykinn/mify/internal/mify/workspace"
	"github.com/chebykinn/mify/pkg/generator"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

func CreateWorkspace(ctx *core.Context, basePath string, name string) error {
	pool, err := util.NewJobPool(ctx, config.GetCacheDirectory(filepath.Join(basePath, name)), runtime.NumCPU())
	if err != nil {
		return err
	}
	defer pool.ClosePool()

	pool.AddJob(util.Job{
		Name: "workspace:" + name,
		Func: func(c *core.Context) error {
			return workspace.CreateWorkspace(c, basePath, name)
		},
	})

	return handleError(pool, pool.Run())
}

func CreateService(ctx *core.Context, workspacePath string, language string, name string) error {
	workspaceContext, pool, err := makeCmdContext(ctx, workspacePath)
	if err != nil {
		return err
	}
	defer pool.ClosePool()

	pool.AddJob(util.Job{
		Name: "create:" + name,
		Func: func(c *core.Context) error {
			return service.CreateService(c, workspaceContext, mifyconfig.ServiceLanguage(language), name)
		},
	})
	if err := pool.Run(); err != nil {
		return handleError(pool, err)
	}
	return handleError(pool, service.Generate(ctx, pool, workspaceContext, name))
}

func CreateFrontend(ctx *core.Context, workspacePath string, template string, name string) error {
	workspaceContext, pool, err := makeCmdContext(ctx, workspacePath)
	if err != nil {
		return err
	}
	defer pool.ClosePool()

	pool.AddJob(util.Job{
		Name: "create:" + service.ApiGatewayName,
		Func: func(c *core.Context) error {
			_, err := service.TryCreateApiGateway(c, workspaceContext)
			return err
		},
	})
	if err := pool.Run(); err != nil {
		return handleError(pool, err)
	}
	// TODO: skip generation if api gateway already exists
	err = handleError(pool, service.Generate(ctx, pool, workspaceContext, service.ApiGatewayName))
	if err != nil {
		return err
	}

	pool.AddJob(util.Job{
		Name: "create:" + name,
		Func: func(c *core.Context) error {
			return service.CreateFrontend(c, workspaceContext, template, name)
		},
	})

	if err := pool.Run(); err != nil {
		return handleError(pool, err)
	}

	return handleError(pool, service.Generate(ctx, pool, workspaceContext, name))
}

func AddClient(ctx *core.Context, workspacePath string, name string, clientName string) error {
	workspaceContext, pool, err := makeCmdContext(ctx, workspacePath)
	if err != nil {
		return err
	}
	defer pool.ClosePool()

	err = client.AddClient(ctx, workspaceContext, name, clientName)
	if err != nil {
		return err
	}

	return handleError(pool, service.Generate(ctx, pool, workspaceContext, name))
}

func RemoveClient(ctx *core.Context, workspacePath string, name string, clientName string) error {
	workspaceContext, pool, err := makeCmdContext(ctx, workspacePath)
	if err != nil {
		return err
	}
	defer pool.ClosePool()

	err = client.RemoveClient(ctx, workspaceContext, name, clientName)
	if err != nil {
		return err
	}

	return handleError(pool, service.Generate(ctx, pool, workspaceContext, name))
}

func ServiceGenerate(ctx *core.Context, workspacePath string, name string) error {
	genPipeline := generator.BuildServicePipeline()
	if err := genPipeline.Execute(); err != nil {
		return err
	}

	workspaceContext, pool, err := makeCmdContext(ctx, workspacePath)
	if err != nil {
		return err
	}
	defer pool.ClosePool()

	return handleError(pool, service.Generate(ctx, pool, workspaceContext, name))
}

func Cleanup() error {
	err := docker.Cleanup(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handleError(pool *util.JobPool, err error) error {
	var jerr util.JobError
	if errors.As(err, &jerr) {
		util.ShowJobError(pool, jerr)
		return jerr.Err
	}
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func makeCmdContext(ctx *core.Context, workspacePath string) (workspace.Context, *util.JobPool, error) {
	workspaceContext, err := workspace.InitContext(workspacePath)
	if err != nil {
		return workspace.Context{}, nil, err
	}

	pool, err := util.NewJobPool(ctx, config.GetCacheDirectory(workspaceContext.BasePath), runtime.NumCPU())
	if err != nil {
		return workspace.Context{}, nil, err
	}
	return workspaceContext, pool, nil
}
