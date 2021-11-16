package service

import (
	"fmt"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

func CreateService(ctx *core.Context, wspContext workspace.Context, name string) error {
	fmt.Printf("Creating service: %s\n", name)

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

	if err := config.SaveServiceConfig(wspContext.BasePath, name, conf); err != nil {
		return err
	}

	return nil
}
