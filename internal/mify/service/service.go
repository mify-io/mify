package service

import (
	"fmt"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

func CreateService(ctx *core.Context, wspContext workspace.Context, name string) error {
	fmt.Printf("creating service %s\n", name)

	repo := fmt.Sprintf("%s/%s/%s",
		wspContext.Config.GitHost,
		wspContext.Config.GitNamespace,
		wspContext.Config.GitRepository)
	context := Context{
		ServiceName: name,
		Repository:  repo,
		GoModule:    repo + "/go_services",
		Workspace:   wspContext,
	}

	if err := RenderTemplateTree(ctx, context); err != nil {
		return err
	}

	conf := config.ServiceConfig{
		ServiceName: name,
	}

	confPath := filepath.Join(wspContext.BasePath, "go_services/cmd", name)
	if err := config.SaveServiceConfig(confPath, conf); err != nil {
		return err
	}

	return nil
}
