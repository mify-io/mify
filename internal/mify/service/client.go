package service

import (
	"fmt"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/workspace"
)


func AddClient(ctx *core.Context, workspaceContext workspace.Context, name string, clientName string) error {
	fmt.Printf("Adding client: %s to %s\n", name, clientName)
	serviceConf, err := config.ReadServiceConfig(workspaceContext.BasePath, name)
	if err != nil {
		return err
	}

	_, err = config.ReadServiceConfig(workspaceContext.BasePath, clientName)
	if err != nil {
		return err
	}

	if serviceConf.OpenAPI.Clients == nil {
		serviceConf.OpenAPI.Clients = map[string]config.ServiceOpenAPIClientConfig{}
	}
	serviceConf.OpenAPI.Clients[clientName] = config.ServiceOpenAPIClientConfig{}
	err = config.SaveServiceConfig(workspaceContext.BasePath, name, serviceConf)
	if err != nil {
		return err
	}

	return nil
}

func RemoveClient(ctx *core.Context, workspaceContext workspace.Context, name string, clientName string) error {
	fmt.Printf("Removing client: %s to %s\n", name, clientName)
	serviceConf, err := config.ReadServiceConfig(workspaceContext.BasePath, name)
	if err != nil {
		return err
	}

	if serviceConf.OpenAPI.Clients == nil {
		return nil
	}

	delete(serviceConf.OpenAPI.Clients, clientName)
	err = config.SaveServiceConfig(workspaceContext.BasePath, name, serviceConf)
	if err != nil {
		return err
	}

	return nil
}
