package client

import (
	"fmt"

	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/workspace"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

func AddClient(ctx *core.Context, workspaceContext workspace.Context, name string, clientName string) error {
	fmt.Printf("Adding client: %s to %s\n", name, clientName)
	serviceConf, err := mifyconfig.ReadServiceConfig(workspaceContext.BasePath, name)
	if err != nil {
		return err
	}

	_, err = mifyconfig.ReadServiceConfig(workspaceContext.BasePath, clientName)
	if err != nil {
		return err
	}

	if serviceConf.OpenAPI.Clients == nil {
		serviceConf.OpenAPI.Clients = map[string]mifyconfig.ServiceOpenAPIClientConfig{}
	}
	serviceConf.OpenAPI.Clients[clientName] = mifyconfig.ServiceOpenAPIClientConfig{}
	err = mifyconfig.SaveServiceConfig(workspaceContext.BasePath, name, serviceConf)
	if err != nil {
		return err
	}

	return nil
}

func RemoveClient(ctx *core.Context, workspaceContext workspace.Context, name string, clientName string) error {
	fmt.Printf("Removing client: %s to %s\n", name, clientName)
	serviceConf, err := mifyconfig.ReadServiceConfig(workspaceContext.BasePath, name)
	if err != nil {
		return err
	}

	if serviceConf.OpenAPI.Clients == nil {
		return nil
	}

	delete(serviceConf.OpenAPI.Clients, clientName)
	err = mifyconfig.SaveServiceConfig(workspaceContext.BasePath, name, serviceConf)
	if err != nil {
		return err
	}

	return nil
}
