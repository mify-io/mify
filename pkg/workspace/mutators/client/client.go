package client

import (
	"fmt"

	"github.com/chebykinn/mify/pkg/mifyconfig"
	"github.com/chebykinn/mify/pkg/workspace/mutators"
)

func AddClient(mutContext *mutators.MutatorContext, name string, clientName string) error {
	fmt.Printf("Adding client: %s to %s\n", name, clientName)

	serviceConf, err := mifyconfig.ReadServiceConfig(mutContext.GetDescription().BasePath, name)
	if err != nil {
		return err
	}

	_, err = mifyconfig.ReadServiceConfig(mutContext.GetDescription().BasePath, clientName)
	if err != nil {
		return err
	}

	if serviceConf.OpenAPI.Clients == nil {
		serviceConf.OpenAPI.Clients = map[string]mifyconfig.ServiceOpenAPIClientConfig{}
	}
	serviceConf.OpenAPI.Clients[clientName] = mifyconfig.ServiceOpenAPIClientConfig{}
	err = mifyconfig.SaveServiceConfig(mutContext.GetDescription().BasePath, name, serviceConf)
	if err != nil {
		return err
	}

	return nil
}

func RemoveClient(mutContext *mutators.MutatorContext, name string, clientName string) error {
	fmt.Printf("Removing client: %s to %s\n", name, clientName)

	serviceConf, err := mifyconfig.ReadServiceConfig(mutContext.GetDescription().BasePath, name)
	if err != nil {
		return err
	}

	if serviceConf.OpenAPI.Clients == nil {
		return nil
	}

	delete(serviceConf.OpenAPI.Clients, clientName)
	err = mifyconfig.SaveServiceConfig(mutContext.GetDescription().BasePath, name, serviceConf)
	if err != nil {
		return err
	}

	return nil
}
